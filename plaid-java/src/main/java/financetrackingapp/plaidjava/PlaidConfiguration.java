package financetrackingapp.plaidjava;

import com.plaid.client.ApiClient;
import com.plaid.client.request.PlaidApi;
import org.springframework.beans.factory.annotation.Value;
import org.springframework.context.annotation.Bean;
import org.springframework.context.annotation.Configuration;

import java.util.HashMap;
import java.util.Map;

@Configuration
public class PlaidConfiguration {

    @Value("${plaid.client.id}")
    private String plaidClientId;

    @Value("${plaid.secret}")
    private String plaidSecret;

    @Value("${plaid.public.key}")
    private String plaidPublicKey;

    @Value("plaid.environment")
    private String plaidEnv;

    private PlaidApi plaidClient;


    @Bean
    public PlaidApi plaidClient() {
        Map<String, String> apiKeys = new HashMap<>();
        apiKeys.put("clientId", plaidClientId);
        apiKeys.put("secret", plaidSecret);

        ApiClient apiClient = new ApiClient(apiKeys);

        String plaidEnvValue = switch (plaidEnv.toLowerCase()) {
            case "development" -> ApiClient.Development;
            case "production" -> ApiClient.Production;
            default -> ApiClient.Sandbox;
        };
        apiClient.setPlaidAdapter(plaidEnvValue);

        plaidClient = apiClient.createService(PlaidApi.class);

        return plaidClient;

    }
}


