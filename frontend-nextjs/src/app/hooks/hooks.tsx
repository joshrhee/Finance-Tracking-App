import { CREATE_LINK_TOKEN_URL } from "../api/api";
import { useQuery } from "@tanstack/react-query";

export function useCreateLinkToken() {
    const { data, error, isFetching, refetch } = useQuery({
        queryKey: ["linkToken"],
        queryFn: async () =>
            await fetch(CREATE_LINK_TOKEN_URL, {
                method: "POST",
                headers: {
                    "Content-Type": "application/json"
                },
                body: JSON.stringify({
                    clientUserId: "testUserId"
                })
            })
                .then(async (result) => await result.json())
                .then((jsonData) => jsonData)
                .catch((error) => {
                    console.error("Link Token error", error);
                }),
        enabled: false
    });

    return { data, error, isFetching, refetch };
}
