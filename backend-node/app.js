// API Routes
// public token: https://{{env_url}}/sandbox/public_token/create
// access token: https://{{env_url}}/item/public_token/exchange
// Balance: https://{{env_url}}/accounts/balance/get
const { v4 } = require("uuid");

const {
    Configuration,
    PlaidApi,
    Products,
    PlaidEnvironments,
    transac,
} = require("plaid");

const express = require("express");
const dotenv = require("dotenv");
const bodyParser = require("body-parser");
const cors = require("cors");

dotenv.config();

const port = process.env.APP_PORT || 8080;
const PLAID_CLIENT_ID = process.env.PLAID_CLIENT_ID;
const PLAID_SECRET = process.env.PLAID_SECRET;
const PLAID_ENV = process.env.PLAID_ENV || "sandbox";

// PLAID_PRODUCTS is a comma-separated list of products to use when initializing
// Link. Note that this list must contain 'assets' in order for the app to be
// able to create and retrieve asset reports.
const PLAID_PRODUCTS = (
    process.env.PLAID_PRODUCTS || Products.Transactions
).split(",");

// PLAID_COUNTRY_CODES is a comma-separated list of countries for which users
// will be able to select institutions from.
const PLAID_COUNTRY_CODES = (process.env.PLAID_COUNTRY_CODES || "US").split(
    ","
);

// Parameter used for OAuth in Android. This should be the package name of your app,
// e.g. com.plaid.linksample
const PLAID_ANDROID_PACKAGE_NAME = process.env.PLAID_ANDROID_PACKAGE_NAME || "";

// Parameters used for the OAuth redirect Link flow.
//
// Set PLAID_REDIRECT_URI to 'http://localhost:3000'
// The OAuth redirect flow requires an endpoint on the developer's website
// that the bank website should redirect to. You will need to configure
// this redirect URI for your client ID through the Plaid developer dashboard
// at https://dashboard.plaid.com/team/api.
const PLAID_REDIRECT_URI = process.env.PLAID_REDIRECT_URI || "";

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

// Initialize the Plaid client
// Find your API keys in the Dashboard (https://dashboard.plaid.com/account/keys)

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

// Sandbox에서 주는 토큰으로 일단하지만, 실제 프로덕션에서는 진짜로 로그인해야함!!!!!
PUBLIC_TOKEN = process.env.PUBLIC_TOKEN;
ACCESS_TOKEN = process.env.ACCESS_TOKEN;
const CLIENT_USER_ID = v4(); // 이 uuid는 나만씀, 결국 다른 사람들도 이거 쓰려면 로그인 정보가지고 uuid만들어야함!!!

const institutionIDs = {
    chase: "ins_270",
};

let institutions = [];

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

const app = express();
app.use(bodyParser.json()); // for parsing application/json
app.use(bodyParser.urlencoded({ extended: true }));
app.use(cors());

app.get("/", (req, res) => {
    res.send("Hello world1111");
});

app.post("/create_user", function (request, response, next) {
    Promise.resolve()
        .then(async function () {})
        .catch(next);
});

// Create a link token with configs which we can then use to initialize Plaid Link client-side.
// See https://plaid.com/docs/#create-link-token
app.post("/create_link_token", function (request, response, next) {
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
            // console.log(
            //     "after create link token, createTokenResponse.data: ",
            //     createTokenResponse.data
            // );
            // LINK_TOKEN = createTokenResponse.data.link_token;
            response.json(createTokenResponse.data);
        })
        .catch(next);
});

// // 일단 쓰지말기!!!!!!
// app.post('/create_public_token', function (request, response, next) {
//   Promise.resolve()
//     .then(async function () {
//       const publicTokenRequest = {
//           institution_id: "insins_03",
//           initial_products: PLAID_PRODUCTS,
//         };
//       console.log("institution_id: ", publicTokenRequest.institution_id)
//       try {
//           const publicTokenResponse = await client.sandboxPublicTokenCreate(
//               publicTokenRequest,
//           );
//           PUBLIC_TOKEN = publicTokenResponse.data.public_token;
//           console.log("Public token: ", PUBLIC_TOKEN);
//           // The generated public_token can now be exchanged
//           // for an access_token
//           const exchangeRequest = {
//               public_token: PUBLIC_TOKEN,
//           };
//           const exchangeTokenResponse = await client.itemPublicTokenExchange(
//               exchangeRequest,
//           );
//           ACCESS_TOKEN = exchangeTokenResponse.data.access_token;
//           console.log("Access token: ", ACCESS_TOKEN);
//       } catch (err) {
//           console.log("Generating public token or access tokein is failed");
//           console.log(err.response);
//       }
//     })
//     .catch(next);
// });

// app.post('/create_link_token', function (request, response, next) {
//   Promise.resolve()
//     .then(async function () {
//       const {client_user_id} = request.body;

//       console.log("Second link token is used!!!!")

//       const configs = {
//         user: {
//           // This should correspond to a unique id for the current user.
//           client_user_id: client_user_id,
//         },
//         client_name: 'Sang June Rhee',
//         products: PLAID_PRODUCTS,
//         country_codes: PLAID_COUNTRY_CODES,
//         language: 'en',
//       };

//       if (PLAID_REDIRECT_URI !== '') {
//         configs.redirect_uri = PLAID_REDIRECT_URI;
//       }

//       if (PLAID_ANDROID_PACKAGE_NAME !== '') {
//         configs.android_package_name = PLAID_ANDROID_PACKAGE_NAME;
//       }
//       const createTokenResponse = await client.linkTokenCreate(configs);
//       // prettyPrintResponse(createTokenResponse);
//       response.json(createTokenResponse.data);
//     })
//     .catch(next);
// })

app.post("/get_access_token", function (request, response, next) {
    Promise.resolve()
        .then(async function () {
            const { public_token } = request.body;
            console.log("public token: ", public_token);

            const tokenResponse = await client
                .itemPublicTokenExchange({
                    public_token: public_token,
                })
                .then((res) => {
                    ACCESS_TOKEN = res.data.access_token;
                    ITEM_ID = res.data.item_id;

                    console.log("GET_ACCESS_TOKEN SUCCEED!!!!!");
                    // console.log("Access token: ", ACCESS_TOKEN);
                    // console.log("ITEM_ID: ", ITEM_ID);

                    // return response.sendStatus(200);
                })
                .catch((err) => {
                    console.log("GET_ACCESS_TOKEN FAILED");
                    // console.log("error: ", err);
                })

                .then(async function () {
                    const transactionRequest = {
                        access_token: ACCESS_TOKEN,
                        start_date: FirstDayOfPreviousMonth,
                        end_date: LastDayOfPreviousMonth,
                    };

                    console.log(
                        "FirstDayOfPreviousMonth: ",
                        FirstDayOfPreviousMonth
                    );
                    console.log(
                        "LastDayOfPreviousMonth: ",
                        LastDayOfPreviousMonth
                    );

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
                    return response.json(editedTransactions);
                    // response.sendStatus(200);
                })
                .catch((err) => {
                    console.log("Getting Transaction FAILED");
                    // console.log("error: ", err);
                });
        })
        .catch(next);
});

// // Retrieve Transactions for an Item
// // https://plaid.com/docs/#transactions
// app.get('/transactions_sync', function (request, response, next) {
//     Promise.resolve()
//     .then(async function (request, response, next) {
//         // Set cursor to empty to receive all historical updates
//         let cursor = null;

//         const config = {
//             access_token: ACCESS_TOKEN,
//             cursor: cursor,
//         };
//         const transactions = await client.transactionsSync(config);
//         const data = transactions.data;
//         console.log("transaction api call data: ", data);
//     })
//     .catch(next);
//   });

app.get("/transactions_get", function (request, response, next) {
    Promise.resolve()
        .then(async function () {
            const transactionRequest = {
                access_token: ACCESS_TOKEN,
                start_date: "2021-01-01",
                end_date: "2021-05-10",
            };

            const client = new PlaidApi(configuration);
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
                    start_date: "2021-01-01",
                    end_date: "2021-05-10",
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
            response.json(editedTransactions);
            // response.sendStatus(200);
        })
        .catch(next);
});

// Retrieve real-time Balances for each of an Item's accounts
// https://plaid.com/docs/#balance
app.get("/balance", function (request, response, next) {
    Promise.resolve()
        .then(async function () {
            const client = new PlaidApi(configuration);
            const balanceResponse = await client.accountsBalanceGet({
                access_token: ACCESS_TOKEN,
            });
            console.log(
                "Balance api call, balanceResponse.data: ",
                balanceResponse.data
            );
            // prettyPrintResponse(balanceResponse);
            response.json(balanceResponse.data);
        })
        .catch(next);
});

app.listen(port, () => {
    console.log(`Listending on port ${port}`);
});
