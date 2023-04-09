package financetrackingapp.plaidjava.controller;

import com.plaid.client.model.Transaction;
import com.plaid.client.model.TransactionsGetRequest;
import com.plaid.client.model.TransactionsGetResponse;
import com.plaid.client.request.PlaidApi;
import financetrackingapp.plaidjava.domain.CustomTransaction;
import financetrackingapp.plaidjava.services.PlaidAuthService;

import lombok.RequiredArgsConstructor;
import org.springframework.http.HttpStatus;
import org.springframework.http.ResponseEntity;
import org.springframework.stereotype.Controller;
import org.springframework.web.bind.annotation.GetMapping;
import retrofit2.Response;

import java.text.SimpleDateFormat;
import java.time.LocalDate;
import java.time.ZoneId;
import java.util.*;


@Controller
@RequiredArgsConstructor
public class TransactionController {
    private final PlaidApi plaidClient;
    private final PlaidAuthService authService;

    @GetMapping("/get_transactions")
    public ResponseEntity<Object> getTransactions(
    ) {
        try {
            SimpleDateFormat simpleDateFormat = new SimpleDateFormat("yyyy-MM-dd");

            Calendar calendar = Calendar.getInstance();
            calendar.add(Calendar.DATE, -30);
            LocalDate startDate = calendar.getTime().toInstant().atZone(ZoneId.systemDefault()).toLocalDate();
            LocalDate endDate = new Date().toInstant().atZone(ZoneId.systemDefault()).toLocalDate();

            TransactionsGetRequest request = new TransactionsGetRequest()
                    .accessToken(this.authService.getAccessToken())
                    .startDate(startDate)
                    .endDate(endDate);
            Response<TransactionsGetResponse> response = plaidClient.transactionsGet(request).execute();

            List<Transaction> transactions = new ArrayList<>();
            transactions.addAll(response.body().getTransactions());

            List<CustomTransaction> customTransactions = new ArrayList<>();
            for (Transaction transaction: transactions) {
                String date = transaction.getDate().toString();
                Double amount = transaction.getAmount();
                List<String> category = transaction.getCategory();
                String name = transaction.getName();

                customTransactions.add(new CustomTransaction(date, amount, category, name));
            }

            return ResponseEntity.ok(customTransactions);


        } catch (Exception e)  {
            return ResponseEntity.status(HttpStatus.INTERNAL_SERVER_ERROR).body(e.getMessage());
        }


    }
}
