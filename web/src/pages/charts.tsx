import { Button } from "@/components/ui/button"
import { ChartConfig, ChartContainer } from "@/components/ui/chart"
import { Dialog, DialogContent, DialogDescription, DialogHeader, DialogTitle, DialogTrigger } from "@/components/ui/dialog"
import { Form, FormControl, FormField, FormItem, FormLabel, FormMessage } from "@/components/ui/form"
import { Select, SelectContent, SelectItem, SelectTrigger, SelectValue } from "@/components/ui/select"
import { Textarea } from "@/components/ui/textarea"
import { useQuery, useQueryClient } from "@tanstack/react-query"
import { Database, Loader2, Plus } from "lucide-react"
import { Bar, BarChart, CartesianGrid, Tooltip, XAxis } from "recharts"

import { zodResolver } from "@hookform/resolvers/zod"
import { useForm } from "react-hook-form"
import { z } from "zod"
import { useState } from "react"
import { toast } from "@/hooks/use-toast"

interface ChartResponse {
	id: string
	type: number
	title: string
	query: string
	databaseId: string
	data: any[]
}

interface DatabaseResponse {
	id: string
	type: number
	name: string
}

export function ChartsPage() {
	const queryClient = useQueryClient()

	const { data: charts, isLoading, error } = useQuery<ChartResponse[]>({
		queryKey: ['charts'],
		queryFn: () =>
			fetch('http://localhost:8080/v1/charts').then(res => res.json())
	})

	const { data: databases } = useQuery<DatabaseResponse[]>({
		queryKey: ['databases'],
		queryFn: () =>
			fetch('http://localhost:8080/v1/databases').then(res => res.json())
	})

	const [isSubmitting, setIsSubmitting] = useState(false)
	const [isDialogOpen, setIsDialogOpen] = useState(false)
	const formSchema = z.object({
		databaseId: z.string().
			nonempty({ message: "This field is required" }).
			uuid({ message: "Invalid Database" }),
		prompt: z.string().
			min(2, { message: "Username must be at least 2 characters." }),
	})

	const form = useForm<z.infer<typeof formSchema>>({
		resolver: zodResolver(formSchema),
		defaultValues: {
			databaseId: undefined,
			prompt: ""
		},
	})

	async function onSubmit(values: z.infer<typeof formSchema>) {
		setIsSubmitting(true)
		try {
			const response = await fetch('http://localhost:8080/v1/charts', {
				method: 'POST',
				headers: {
					'Content-Type': 'application/json',
				},
				body: JSON.stringify(values),
			});

			if (!response.ok) {
				throw new Error('Network response was not ok');
			}

			const data = await response.json();
			console.log('Chart created:', data);
			setIsDialogOpen(false)
			form.reset() // Reset form
			toast({
				title: 'Chart created',
				description: 'Chart created successfully',
				variant: 'default',
			});
			await queryClient.invalidateQueries({ queryKey: ['charts'] })
		} catch (error) {
			toast({
				title: 'Error creating chart',
				description: error instanceof Error ? error.message : 'Unknown error',
				variant: 'destructive',
			});
		} finally {
			setIsSubmitting(false)
		}
	}

	if (isLoading) return <div>Loading...</div>
	if (error) return <div>Error loading charts</div>

	return (
		<div className="flex flex-col gap-2">
			<h1 className="scroll-m-20 text-3xl font-semibold tracking-tight first:mt-4">
				Charts
			</h1>
			<div>
				<Dialog open={isDialogOpen} onOpenChange={setIsDialogOpen}>
					<DialogTrigger asChild>
						<Button onClick={() => setIsDialogOpen(true)}>
							<Plus /> Add
						</Button>
					</DialogTrigger>
					<DialogContent className="sm:max-w-[700px]">
						<DialogHeader>
							<DialogTitle>Create a new Chart</DialogTitle>
							<DialogDescription>
								Create a Chart using the Power of AI ✨
							</DialogDescription>
						</DialogHeader>
						<Form {...form}>
							<form onSubmit={form.handleSubmit(onSubmit)} className="space-y-8">
								<FormField
									control={form.control}
									name="databaseId"
									render={({ field }) => (
										<FormItem>
											<FormLabel>Database</FormLabel>
											<FormControl>
												<Select
													onValueChange={field.onChange}
													defaultValue={field.value}
													{...field}>
													<SelectTrigger className="col-span-3">
														<SelectValue placeholder="Select a Database" />
													</SelectTrigger>
													<SelectContent>
														{databases?.map((database) =>
															<SelectItem value={database.id}><div className="flex gap-2 content-center"><Database className="h-4 w-4 m-0.5" /> {database.name}</div></SelectItem>
														)}
													</SelectContent>
												</Select>
											</FormControl>
											<FormMessage />
										</FormItem>
									)}
								/>
								<FormField
									control={form.control}
									name="prompt"
									render={({ field }) => (
										<FormItem>
											<FormLabel>Prompt</FormLabel>
											<FormControl>
												<Textarea {...field} placeholder="e.g.: Create a chart of the total revenue grouped by stores name" className="col-span-3 min-h-40" />
											</FormControl>
											<FormMessage />
										</FormItem>
									)}
								/>
								<Button disabled={isSubmitting} type="submit">
									{isSubmitting ? (
										<>
											<span className="animate-spin mr-2"><Loader2 /></span> Creating...
										</>
									) : (
										'Submit'
									)}
								</Button>
							</form>
						</Form>
					</DialogContent>
				</Dialog>
			</div>
			<div className="grid mt-1 xs:grid-cols-1 sm:grid-cols-2 md:grid-cols-3 gap-4">
				{charts?.map((chart: ChartResponse, i: number) => (
					<div
						key={chart.id}
						className="p-1 rounded-md border aspect-square flex flex-col justify-between"
					>
						<div className="text-center font-normal text-sm">{chart.title}</div>
						<ChartContainer
							config={chartConfig}
							className="aspect-auto h-[250px] w-full"
						>
							<BarChart
								accessibilityLayer
								data={chart.data.sort((a, b) => a.Value > b.Value ? -1 : 1)}
								margin={{
									left: 12,
									right: 12,
								}}
							>
								<CartesianGrid vertical={false} />
								<XAxis
									dataKey="Category"
									tickLine={false}
									axisLine={false}
									tickMargin={8}
									minTickGap={32}
								/>
								<Tooltip />
								<Bar dataKey={"Value"} fill={`hsl(var(--chart-${i % 5 + 1}))`} />
							</BarChart>
						</ChartContainer>
					</div>
				))}
			</div>
		</div>
	)
}

const chartConfig = {
	views: {
		label: "Page Views",
	},
	desktop: {
		label: "Desktop",
		color: "hsl(var(--chart-1))",
	},
	mobile: {
		label: "Mobile",
		color: "hsl(var(--chart-2))",
	},
} satisfies ChartConfig
