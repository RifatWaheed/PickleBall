import React, { createContext, useContext, useEffect, useMemo, useState } from 'react';
import { useQuery, useQueryClient } from '@tanstack/react-query';
import { login as loginApi, logout as logoutApi, me as meApi, refresh as refreshApi, MeResponse } from '@/services/authApi';
import { setAccessTokenMemory } from '@/services/http';

type AuthState = {
    user: MeResponse | null;
    isAuthed: boolean;
    isBootstrapping: boolean;
    login: (email: string, password: string) => Promise<void>;
    logout: () => Promise<void>;
};

const AuthContext = createContext<AuthState | null>(null);

export function useAuth() {
    const ctx = useContext(AuthContext);
    if (!ctx) throw new Error('useAuth must be used inside AuthProvider');
    return ctx;
}

export default function AuthProvider({ children }: { children: React.ReactNode }) {
    const qc = useQueryClient();
    const [bootstrapped, setBootstrapped] = useState(false);

    // 1) On app start: try refresh -> set access token -> fetch me
    const meQuery = useQuery({
        queryKey: ['me'],
        queryFn: meApi,
        enabled: bootstrapped, // only after refresh attempt
        retry: false,
    });

    useEffect(() => {
        (async () => {
            try {
                const r = await refreshApi();
                setAccessTokenMemory(r.accessToken);
            } catch {
                setAccessTokenMemory(null);
            } finally {
                setBootstrapped(true);
            }
        })();
    }, []);

    const login = async (email: string, password: string) => {
        const res = await loginApi({ email, password });
        setAccessTokenMemory(res.accessToken);

        await qc.fetchQuery({
            queryKey: ['me'],
            queryFn: meApi,
        });
    };


    const logout = async () => {
        try {
            await logoutApi();
        } finally {
            setAccessTokenMemory(null);
            qc.clear();
        }
    };

    const value = useMemo<AuthState>(() => {
        const user = meQuery.data ?? null;
        return {
            user,
            isAuthed: !!user,
            isBootstrapping: !bootstrapped || meQuery.isLoading,
            login,
            logout,
        };
    }, [bootstrapped, meQuery.data, meQuery.isLoading]);

    return <AuthContext.Provider value={value}>{children}</AuthContext.Provider>;
}
