import type { Incident, CreateIncidentRequest, Severity } from "./types.js";

interface ApiError {
    error: {
        code: string;
        message: string;
    }
}

export async function isAuthenticated() : Promise<boolean> {
    const res = await fetch("/api/auth/isauthenticated", {
        credentials: "include"
    });
    return res.ok;   
}

async function request<T>(url: string, init?: RequestInit) : Promise<T> {
    const res = await fetch(url, {
        credentials: "include",
        ...init,
    })
    if (res.status == 204) {
        return undefined as T
    }
    
    const data: unknown = await res.json()
    if (!res.ok) {
        const err = data as ApiError
        const errorCode = err.error.code;
        const message = err.error.message;
        throw new Error(`${errorCode}: ${message}`);
    }
    
    return data as T
}

export async function login(username:string, password:string) : Promise<void> {
    return await request<void>("/login", {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify({ username, password }),
    })
}

export async function loadIncidents() : Promise<Incident[]> {
    return await request<Incident[]>("/api/incidents", undefined)
}

export async function logout() : Promise<void> {
    return await request<void>("/api/auth/logout", {
        method: "POST",
    })
}

export async function createIncident(input:CreateIncidentRequest) : Promise<Incident> {
    return await request<Incident>("/api/incidents", {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify(input),
    })
}
