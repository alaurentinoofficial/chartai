import { PanelLeftOpen, PanelLeftClose } from "lucide-react"
import { useSidebar } from "@/components/ui/sidebar"
import { Button } from "@/components/ui/button"

export function Navbar() {
	const {
		open,
		toggleSidebar,
	  } = useSidebar()
	
	
	return (
		<header className="sticky p-3 top-0 z-50 w-full bg-background/95 backdrop-blur supports-[backdrop-filter]:bg-background/60">
			<div className="container flex h-14 items-center">
				<Button variant="ghost" size="icon" onClick={toggleSidebar}>
					{open ? <PanelLeftClose /> : <PanelLeftOpen />}
				</Button>
			</div>
		</header>
	)
}
