package financetrackingapp.plaidjava.controller;

import com.plaid.client.model.ItemPublicTokenExchangeRequest;
import com.plaid.client.model.ItemPublicTokenExchangeResponse;
import com.plaid.client.request.PlaidApi;
import financetrackingapp.plaidjava.services.PlaidAuthService;
import lombok.RequiredArgsConstructor;
import org.springframework.core.env.Environment;
import org.springframework.http.HttpStatus;
import org.springframework.http.ResponseEntity;
import org.springframework.stereotype.Controller;
import org.springframework.web.bind.annotation.CrossOrigin;
import org.springframework.web.bind.annotation.PostMapping;
import org.springframework.web.bind.annotation.RequestParam;
import org.springframework.web.server.ResponseStatusException;
import retrofit2.Response;

@Controller
@RequiredArgsConstructor
public class AccessTokenController {
    private final Environment env;
    private final PlaidApi plaidClient;
    private final PlaidAuthService authService;
    private final TransactionController transactionController;

    @CrossOrigin
    @PostMapping("/get_access_token")
    public ResponseEntity<Object> getAccessToken(@RequestParam("public_token") String publicToken) {
        try {

            ItemPublicTokenExchangeRequest request = new ItemPublicTokenExchangeRequest()
                    .publicToken(publicToken);

            Response<ItemPublicTokenExchangeResponse> response = plaidClient
                    .itemPublicTokenExchange(request)
                    .execute();

            // // plaid 2.2.0 version
//            Response<ItemPublicTokenExchangeResponse> response = this.plaidClient.service()
//                    .itemPublicTokenExchange(new ItemPublicTokenExchangeRequest(publicToken))
//                    .execute();

            this.authService.setAccessToken(response.body().getAccessToken());
            this.authService.setItemId(response.body().getItemId());

            if (authService.getAccessToken() == null) {
                throw new ResponseStatusException(HttpStatus.UNAUTHORIZED, "Access Token is null!!!!");
            }

//            System.out.println("Comment!!!!!!!!, public token: " + publicToken);
//            System.out.println("Comment!!!!!!!!, access token: " + this.authService.getAccessToken());
//            System.out.println("Comment!!!!!!!!, itemId: " + this.authService.getItemId());

            ResponseEntity<Object> transactionResponse = this.transactionController.getTransactions();

            return ResponseEntity.ok(transactionResponse.getBody());

        } catch (Exception e) {
            e.printStackTrace();
            return ResponseEntity.status(HttpStatus.INTERNAL_SERVER_ERROR).body(e.getMessage());
        }
    }
}
