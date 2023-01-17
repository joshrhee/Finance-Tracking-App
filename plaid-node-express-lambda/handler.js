const serverless = require("serverless-http");
const express = require("express");
const cors = require("cors");
const app = express();
app.use(cors());

const {
    Configuration,
    PlaidApi,
    Products,
    PlaidEnvironments,
    transac,
} = require("plaid");
const dotenv = require("dotenv");
dotenv.config();
const bodyParser = require("body-parser");
const { v4: uuidv4 } = require("uuid");

const CLIENT_USER_ID = uuidv4();

const port = process.env.APP_PORT || 8080;
const PLAID_CLIENT_ID = process.env.PLAID_CLIENT_ID;
const PLAID_SECRET = process.env.PLAID_SECRET;
const PLAID_ENV = process.env.PLAID_ENV;
const PLAID_PRODUCTS = ["auth", "transactions", "identity"];
const PLAID_COUNTRY_CODES = ["US", "CA"];

const PLAID_ANDROID_PACKAGE_NAME = "";
const PLAID_REDIRECT_URI = "";
const GIN_MODE = process.env.GIN_MODE;

console.log("PLAID_CLIENT_ID: ", PLAID_CLIENT_ID);

// We store the access_token in memory - in production, store it in a secure
// persistent data store
let ACCESS_TOKEN = null;
let PUBLIC_TOKEN = null;
let ITEM_ID = null;
// The payment_id is only relevant for the UK/EU Payment Initiation product.
// We store the payment_id in memory - in production, store it in a secure
// persistent data store along with the Payment metadata, such as userId .
let PAYMENT_ID = null;
// The transfer_id is only relevant for Transfer ACH product.
// We store the transfer_id in memory - in production, store it in a secure
// persistent data store
let TRANSFER_ID = null;

let LINK_TOKEN = null;

const configuration = new Configuration({
    basePath: PlaidEnvironments[PLAID_ENV],
    baseOptions: {
        headers: {
            "PLAID-CLIENT-ID": PLAID_CLIENT_ID,
            "PLAID-SECRET": PLAID_SECRET,
            "Plaid-Version": "2020-09-14",
        },
    },
});
const client = new PlaidApi(configuration);

function getTransactionDateRange() {
    const date = new Date();
    const firstDay = new Date(date.getFullYear(), date.getMonth() - 1, 1)
        .toISOString()
        .split("T")[0];
    const lastDay = new Date(date.getFullYear(), date.getMonth(), 0)
        .toISOString()
        .split("T")[0];
    return {
        firstDay: firstDay,
        lastDay: lastDay,
    };
}

const PreviousMonthDateRange = getTransactionDateRange();
const FirstDayOfPreviousMonth = PreviousMonthDateRange.firstDay;
const LastDayOfPreviousMonth = PreviousMonthDateRange.lastDay;

app.post("/", (req, res, next) => {
    return res.status(200).json({
        message: "Hello from root!",
    });
});

app.post("/create_link_token", (request, response, next) => {
    Promise.resolve()
        .then(async function () {
            const configs = {
                user: {
                    // This should correspond to a unique id for the current user.
                    client_user_id: CLIENT_USER_ID,
                },
                client_name: "Plaid Quickstart",
                products: PLAID_PRODUCTS,
                country_codes: PLAID_COUNTRY_CODES,
                language: "en",
            };

            if (PLAID_REDIRECT_URI !== "") {
                configs.redirect_uri = PLAID_REDIRECT_URI;
            }

            if (PLAID_ANDROID_PACKAGE_NAME !== "") {
                configs.android_package_name = PLAID_ANDROID_PACKAGE_NAME;
            }
            const createTokenResponse = await client.linkTokenCreate(configs);
            // prettyPrintResponse(createTokenResponse);
            console.log(
                "after create link token, createTokenResponse.data: ",
                createTokenResponse.data
            );
            // LINK_TOKEN = createTokenResponse.data.link_token;
            return response.status(200).json(createTokenResponse.data);
        })
        .catch(next);

    // return res.status(200).json({
    //     message: "Hello from create_link_token path! cleiint: ",
    //     client,
    // });
});

app.post("/get_access_token", function (request, response, next) {
    Promise.resolve()
        .then(async function () {
            const requestBody = JSON.parse(request.apiGateway.event.body);
            const public_token = requestBody.public_token;
            console.log("request.body: ", requestBody);

            console.log("public_token: ", public_token);

            const client = new PlaidApi(configuration);
            const tokenResponse = await client
                .itemPublicTokenExchange({
                    public_token: public_token,
                })
                .then((res) => {
                    console.log("res.data: ", res.data);
                    ACCESS_TOKEN = res.data.access_token;
                    ITEM_ID = res.data.item_id;

                    console.log("GET_ACCESS_TOKEN SUCCEED!!!!!");
                    // console.log("Access token: ", ACCESS_TOKEN);
                    // console.log("ITEM_ID: ", ITEM_ID);

                    // response.sendStatus(200);
                })
                .catch((err) => {
                    console.log("GET_ACCESS_TOKEN FAILED");
                    return response.status(500).json(err.message);
                })

                .then(async function () {
                    const transactionRequest = {
                        access_token: ACCESS_TOKEN,
                        start_date: FirstDayOfPreviousMonth,
                        end_date: LastDayOfPreviousMonth,
                    };

                    const transactionsResponse = await client.transactionsGet(
                        transactionRequest
                    );
                    let transactions = transactionsResponse.data.transactions;
                    const total_transactions =
                        transactionsResponse.data.total_transactions;
                    // Manipulate the offset parameter to paginate
                    // transactions and retrieve all available data
                    while (transactions.length < total_transactions) {
                        const paginatedRequest = {
                            access_token: ACCESS_TOKEN,
                            start_date: FirstDayOfPreviousMonth,
                            end_date: LastDayOfPreviousMonth,
                            options: {
                                offset: transactions.length,
                            },
                        };
                        const paginatedResponse = await client.transactionsGet(
                            paginatedRequest
                        );
                        transactions = transactions.concat(
                            paginatedResponse.data.transactions
                        );
                    }

                    let editedTransactions = [];
                    transactions.map((transaction) => {
                        editedTransactions.push({
                            date: transaction.date,
                            amount: transaction.amount,
                            category: transaction.category,
                            name: transaction.name,
                        });
                    });
                    // console.log("transactions get response, transactions: ", transactions)
                    console.log(
                        "transactions get response, Edited transactions: ",
                        editedTransactions
                    );
                    return response.status(200).json(editedTransactions);
                    // response.sendStatus(200);
                })
                .catch((err) => {
                    console.log("Getting Transaction FAILED");
                    return response.status(500).json(err.message);
                });
        })
        .catch(next);
});

// app.post("/transactions_get", (request, response, next) => {
//     Promise.resolve()
//         .then(async function () {
//             const transactionRequest = {
//                 access_token: ACCESS_TOKEN,
//                 start_date: "2021-01-01",
//                 end_date: "2021-05-10",
//             };

//             const transactionsResponse = await client.transactionsGet(
//                 transactionRequest
//             );
//             let transactions = transactionsResponse.data.transactions;
//             const total_transactions =
//                 transactionsResponse.data.total_transactions;
//             // Manipulate the offset parameter to paginate
//             // transactions and retrieve all available data
//             while (transactions.length < total_transactions) {
//                 const paginatedRequest = {
//                     access_token: ACCESS_TOKEN,
//                     start_date: "2021-01-01",
//                     end_date: "2021-05-10",
//                     options: {
//                         offset: transactions.length,
//                     },
//                 };
//                 const paginatedResponse = await client.transactionsGet(
//                     paginatedRequest
//                 );
//                 transactions = transactions.concat(
//                     paginatedResponse.data.transactions
//                 );
//             }

//             let editedTransactions = [];
//             transactions.map((transaction) => {
//                 editedTransactions.push({
//                     date: transaction.date,
//                     amount: transaction.amount,
//                     category: transaction.category,
//                     name: transaction.name,
//                 });
//             });
//             // console.log("transactions get response, transactions: ", transactions)
//             console.log(
//                 "transactions get response, Edited transactions: ",
//                 editedTransactions
//             );
//             response.json(editedTransactions);
//             // response.sendStatus(200);
//         })
//         .catch(next);
// });

app.use((req, res, next) => {
    return res.status(404).json({
        error: "Not Found",
    });
});

app.listen(port, () => {
    console.log(`Listending on port ${port}`);
});

module.exports.handler = serverless(app);
