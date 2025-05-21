export async function fetchDomain(): Promise<string> {
    const res = await fetch("/api/domain");
    if (!res.ok) throw new Error("Failed to fetch domain");
    return await res.json();
}

export async function fetchIp(): Promise<string> {
    const res = await fetch("/api/ip");
    if (!res.ok) throw new Error("Failed to fetch ip");
    return await res.json();
}

export async function fetchPortMap(): Promise<Record<string, string>> {
    const res = await fetch("/api/config");
    if (!res.ok) throw new Error("Failed to fetch port map");
    return await res.json();
}
