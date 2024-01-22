"use client";

import { useState } from "react";
import { QueryClient, QueryClientProvider } from "@tanstack/react-query";
import { ReactQueryDevtools } from "@tanstack/react-query-devtools";

export default function QueryProvider({
    children
}: {
    children: React.ReactNode;
}) {
    const [queryClient] = useState(
        () =>
            new QueryClient({
                defaultOptions: {
                    //   queries: {
                    //     // With SSR, we usually want to set some default staleTime
                    //     // above 0 to avoid refetching immediately on the client
                    //     staleTime: 4 * 1000,
                    //     refetchInterval: 4 * 1000,
                    //   },
                }
            })
    );

    return (
        <QueryClientProvider client={queryClient}>
            <ReactQueryDevtools />
            {children}
        </QueryClientProvider>
    );
}
