import { QueryClient, QueryClientProvider } from "@tanstack/react-query";
import { Routers } from "@/Routers";
import { ThemeProvider } from "@/components/theme-provider";

export function App() {
	const queryClient = new QueryClient()
	return (
		<QueryClientProvider client={queryClient}>
			<ThemeProvider defaultTheme="light" storageKey="vite-ui-theme">
				<Routers />
			</ThemeProvider>
		</QueryClientProvider>
	)
}
