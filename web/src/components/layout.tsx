import { SidebarInset, SidebarProvider, SidebarTrigger } from "@/components/ui/sidebar";
import { AppSidebar } from "@/components/sidebar";
import { Separator } from "@/components/ui/separator";
import { Breadcrumb, BreadcrumbItem, BreadcrumbLink, BreadcrumbList } from "@/components/ui/breadcrumb";

export default function Layout({ children }: { children: React.ReactNode }) {
	const defaultOpen = true

	return (
		<div className="flex w-full">
			<SidebarProvider defaultOpen={defaultOpen}>
				<div className="relative flex min-h-screen w-full">
					<AppSidebar />
					<SidebarInset>
						<div className="flex flex-col flex-1 w-full">
							<header className="flex h-16 shrink-0 items-center gap-2 transition-[width,height] ease-linear group-has-[[data-collapsible=icon]]/sidebar-wrapper:h-12">
								<div className="flex items-center gap-2 px-4">
									<SidebarTrigger className="-ml-1" />
									<Separator orientation="vertical" className="mr-2 h-4" />
									<Breadcrumb>
										<BreadcrumbList>
											<BreadcrumbItem className="hidden md:block">
												<BreadcrumbLink href="#">
													Charts
												</BreadcrumbLink>
											</BreadcrumbItem>
											{ /* <BreadcrumbSeparator className="hidden md:block" /> */ }
										</BreadcrumbList>
									</Breadcrumb>
								</div>
							</header>
							<main className="flex-1 px-4 w-full">
								{children}
							</main>
						</div>
					</SidebarInset>
				</div>
			</SidebarProvider>
		</div>
	)
}
