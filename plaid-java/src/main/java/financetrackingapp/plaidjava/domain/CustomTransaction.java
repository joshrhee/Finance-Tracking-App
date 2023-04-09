package financetrackingapp.plaidjava.domain;

import java.util.List;

public class CustomTransaction implements TransactionInterface {
    private final String date;
    private final Double amount;
    private final List<String> category;
    private final String name;

    public CustomTransaction(String date, Double amount, List<String> category, String name) {
        this.date = date;
        this.amount = amount;
        this.category = category;
        this.name = name;
    }

    @Override
    public String getDate() {
        return date;
    }

    @Override
    public Double getAmount() {
        return amount;
    }

    @Override
    public List<String> getCategory() {
        return category;
    }

    @Override
    public String getName() {
        return name;
    }
}
