"use client";

import { useEffect, useState } from "react";
import {
	Table,
	TableBody,
	TableCell,
	TableHead,
	TableHeader,
	TableRow
} from "@/components/ui/table";

import { fetchDomain, fetchIp, fetchPortMap } from "./api";
import { Button } from "@/components/ui/button";
import { ExternalLink } from "lucide-react";

export default function HomePage() {
	const [domain, setDomain] = useState<string>("");
	const [ip, setIp] = useState<string>("");
	const [portMap, setPortMap] = useState<Record<string, string>>({});
	const [loading, setLoading] = useState(true);
	const [error, setError] = useState<string | null>(null);

	useEffect(() => {
		const loadData = async () => {
			try {
				const [domain, ip, ports] = await Promise.all([
					fetchDomain(),
					fetchIp(),
					fetchPortMap()
				]);
				setDomain(domain);
				setIp(ip);
				setPortMap(ports);
			} catch (err) {
				console.error(err);
				setError("Fail to fetch data");
			} finally {
				setLoading(false);
			}
		};

		loadData();
	}, []);

	if (loading) return <p className="p-4">Loading...</p>;
	if (error) return <p className="p-4 text-red-500">{error}</p>;

	const formatAddress = (ip: string, port: string) => {
		if (port.startsWith(":")) {
			return `${ip}${port}`;
		}
		return port;
	};

	return (
		<Table>
			<TableHeader>
				<TableRow>
					<TableHead>subdomain</TableHead>
					<TableHead>target</TableHead>
					<TableHead></TableHead>
				</TableRow>
			</TableHeader>
			<TableBody>
				{Object.entries(portMap).map(([subdomain, port]) => (
					<TableRow key={subdomain}>
						<TableCell>{subdomain}</TableCell>
						<TableCell>{formatAddress(ip, port)}</TableCell>
						<TableCell>
							<Button
								variant="ghost" className="h-8 w-8 p-0"
								onClick={() => window.open(`https://${subdomain}.${domain}/`, "_blank")}
							>
								<ExternalLink className="h-4 w-4" />
							</Button>
						</TableCell>
					</TableRow>
				))}
			</TableBody>
		</Table >
	);
}
