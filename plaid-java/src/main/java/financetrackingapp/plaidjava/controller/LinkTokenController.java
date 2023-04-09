//package financetrackingapp.plaidjava.controller;
//
//import com.plaid.client.model.*;
//import com.plaid.client.request.PlaidApi;
//import lombok.RequiredArgsConstructor;
//import org.springframework.beans.factory.annotation.Autowired;
//import org.springframework.beans.factory.annotation.Value;
//import org.springframework.http.ResponseEntity;
//import org.springframework.stereotype.Controller;
//import org.springframework.web.bind.annotation.PostMapping;
//import retrofit2.Response;
//
//import java.io.IOException;
//import java.util.ArrayList;
//import java.util.Arrays;
//import java.util.Date;
//import java.util.List;
//
//@Controller
//@RequiredArgsConstructor
//public class LinkTokenController {
//
//    private final PlaidApi plaidClient;
//    private final List<String> plaidProducts;
//    private final List<String> countryCodes;
//    private final String redirectUri;
//    private final List<Products> correctedPlaidProducts;
//    private final List<CountryCode> correctedCountryCodes;
//
//    @Value("${plaid.client.id}")
//    private String clientId;
//
//    @Value("${plaid.secret}")
//    private String secret;
//
//    // Autowired inject PlaidClient bean that is defined in PlaidConfiguration.java
//    @Autowired
//    public LinkTokenController(PlaidApi plaidClient) {
//        this.plaidClient = plaidClient;
//        this.plaidProducts = new ArrayList<>();
//        this.countryCodes = new ArrayList<>();
//        this.redirectUri = "";
//        this.correctedPlaidProducts = new ArrayList<>();
//        this.correctedCountryCodes = new ArrayList<>();
//    }
//
//    @PostMapping("/create_link_token")
//    public ResponseEntity<Object> createLinkToken() throws IOException {
//
//        System.out.println("Comment plaidClient: " + plaidClient);
//
//        String clientUserID = Long.toString((new Date()).getTime());
//        LinkTokenCreateRequestUser user = new LinkTokenCreateRequestUser()
//                .clientUserId("clientUserID");
//
//        System.out.println("Comment user: " + user);
//
//        for (String plaidProduct : this.plaidProducts) {
//            this.correctedPlaidProducts.add(Products.fromValue(plaidProduct));
//        }
//
//        for (String countryCode : this.countryCodes) {
//            this.correctedCountryCodes.add(CountryCode.fromValue(countryCode));
//        }
//
//        LinkTokenCreateRequest request = new LinkTokenCreateRequest()
//                .user(user)
//                .clientName("Quickstart Client")
//                .products(Arrays.asList(Products.AUTH))
//                .countryCodes(Arrays.asList(CountryCode.US))
//                .language("en")
//                .redirectUri(this.redirectUri)
//                        .clientId(clientId)
//                                .secret(secret);
//
//        System.out.println("Comment request: " + request);
//
//        Response<LinkTokenCreateResponse> response = plaidClient
//                .linkTokenCreate(request)
//                .execute();
//
//        System.out.println("comment response.body(): " + response);
//
//        String linkToken = response.body().getLinkToken();
//
//
//
//        return ResponseEntity.ok(linkToken);
//    }
//}
