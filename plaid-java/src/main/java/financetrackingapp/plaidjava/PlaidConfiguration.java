package financetrackingapp.plaidjava;

//import com.fasterxml.jackson.annotation.JsonProperty;
import com.plaid.client.ApiClient;
//import com.plaid.client.PlaidClient;
import com.plaid.client.request.PlaidApi;
//import io.swagger.annotations.Api;
//import org.hibernate.validator.constraints.NotEmpty;
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

        String plaidEnvValue = "";
        switch (plaidEnv.toLowerCase()) {
            case "sandbox":
                plaidEnvValue = ApiClient.Sandbox;
                break;
            case "development":
                plaidEnvValue = ApiClient.Development;
                break;
            case "production":
                plaidEnvValue = ApiClient.Production;
                break;
            default:
                plaidEnvValue = ApiClient.Sandbox;
        }
        apiClient.setPlaidAdapter(plaidEnvValue);

        plaidClient = apiClient.createService(PlaidApi.class);

        return plaidClient;

        // // plaid version 2.2.0
//        PlaidClient.Builder clientBuilder = PlaidClient.newBuilder()
//                .clientIdAndSecret(plaidClientId, plaidSecret)
//                .publicKey(plaidPublicKey);


//        switch (plaidEnv.toLowerCase()) {
//            case "sandbox":
//                clientBuilder.sandboxBaseUrl();
//                break;
//            case "development":
//                clientBuilder.developmentBaseUrl();
//                break;
//            case "production":
//                clientBuilder.productionBaseUrl();
//                break;
//            default:
//                clientBuilder.sandboxBaseUrl();
//        }

//        return clientBuilder.build();
    }
}


//public class PlaidConfiguration extends Configuration {
//    @NotEmpty
//    private String plaidClientID;
//
//    @NotEmpty
//    private String plaidSecret;
//
//    @NotEmpty
//    private String plaidEnv;
//
//    @NotEmpty
//    private String plaidProducts;
//
//    @NotEmpty
//    private String plaidCountryCodes;
//
//    // Parameters used for the OAuth redirect Link flow.
//
//    // Set PLAID_REDIRECT_URI to 'http://localhost:3000/'
//    // The OAuth redirect flow requires an endpoint on the developer's website
//    // that the bank website should redirect to. You will need to configure
//    // this redirect URI for your client ID through the Plaid developer dashboard
//    // at https://dashboard.plaid.com/team/api.
//    private String plaidRedirectUri;
//
//    @JsonProperty
//    public String getPlaidClientID() {
//        return plaidClientID;
//    }
//
//    @JsonProperty
//    public String getPlaidSecret() {
//        return plaidSecret;
//    }
//
//    @JsonProperty
//    public String getPlaidEnv() {
//        return plaidEnv;
//    }
//
//    @JsonProperty
//    public String getPlaidProducts() {
//        return plaidProducts;
//    }
//
//    @JsonProperty
//    public String getPlaidCountryCodes() {
//        return plaidCountryCodes;
//    }
//
//    @JsonProperty
//    public String getPlaidRedirectUri() {
//        return plaidRedirectUri;
//    }
//}
