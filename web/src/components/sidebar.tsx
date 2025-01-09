import { ChartArea, Database, GalleryVerticalEnd, AudioWaveform, Command } from "lucide-react"
import { Sidebar, SidebarContent, SidebarGroup, SidebarGroupContent, SidebarHeader, SidebarMenu, SidebarMenuButton, SidebarMenuItem } from "@/components/ui/sidebar"
import { TeamSwitcher } from "@/components/team-switcher"
import { Link } from "react-router-dom";

const items = [
	// {
	// 	title: "Home",
	// 	url: "/",
	// 	icon: Home,
	// },
	{
		title: "Charts",
		url: "/charts",
		icon: ChartArea,
	},
	{
		title: "Databases",
		url: "/databases",
		icon: Database,
	},
	// {
	// 	title: "Settings",
	// 	url: "/settings",
	// 	icon: Settings,
	// },
]

const teams = [
	{
		name: "Acme Inc",
		logo: GalleryVerticalEnd,
		plan: "Enterprise",
	},
	{
		name: "Acme Corp.",
		logo: AudioWaveform,
		plan: "Startup",
	},
	{
		name: "Evil Corp.",
		logo: Command,
		plan: "Free",
	},
]

export function AppSidebar() {
	return (
		<Sidebar collapsible="icon">
			<SidebarHeader>
				<TeamSwitcher teams={teams} />
			</SidebarHeader>
			<SidebarContent>
				<SidebarGroup>
					<SidebarGroupContent>
						<SidebarMenu>
							{items.map((item) => (
								<SidebarMenuItem key={item.title}>
									<SidebarMenuButton asChild>
										<Link to={item.url}>
											<item.icon />
											<span>{item.title}</span>
										</Link>
									</SidebarMenuButton>
								</SidebarMenuItem>
							))}
						</SidebarMenu>
					</SidebarGroupContent>
				</SidebarGroup>
			</SidebarContent>
		</Sidebar>
	)
}
