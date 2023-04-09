package financetrackingapp.plaidjava.domain;

import java.util.List;

public interface TransactionInterface {
    String getDate();
    Double getAmount();
    List<String> getCategory();
    String getName();
}

