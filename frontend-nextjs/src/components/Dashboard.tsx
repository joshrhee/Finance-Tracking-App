"use client";

import { TableComponent } from "@/components/TableComponent";
import { Button } from "@/components/ui/button";
import {
    HydrationBoundary,
    dehydrate,
    QueryClient
} from "@tanstack/react-query";
import { useCreateLinkToken } from "@/app/hooks/hooks";
import { useEffect } from "react";
import { AlertComponent } from "@/components/AlertComponent";
import { LoadingSpinner } from "./LoadingSpinner";

export function Dashboard() {
    const queryClient = new QueryClient();
    const {
        data: linkToken,
        error,
        isFetching,
        refetch
    } = useCreateLinkToken();

    // // Does it necessary for prefetch?
    // await queryClient.prefetchQuery({
    //     queryKey: ["transactions"],
    //     queryFn: () => {}
    // })

    useEffect(() => {
        console.log("linkToken", linkToken);
    }, [linkToken]);

    return (
        <main className="flex min-h-screen flex-col items-center justify-between p-24">
            <HydrationBoundary state={dehydrate(queryClient)}>
                <Button
                    variant="outline"
                    disabled={isFetching}
                    onClick={() => {
                        console.log("Link Token button is clicked");
                        refetch();
                    }}
                >
                    {isFetching ? (
                        <LoadingSpinner className="spinner" />
                    ) : (
                        "Click this button first (Create Lnk Token)"
                    )}
                </Button>
                {error && (
                    <AlertComponent
                        title="Heads up!!"
                        description={error.message}
                    />
                )}
                <Button
                    variant="outline"
                    // onClick={() => {
                    //     console.log("Transaction button is clicked");
                    // }}
                >
                    Get Transactions information
                </Button>
                <TableComponent />
            </HydrationBoundary>
        </main>
    );
}
