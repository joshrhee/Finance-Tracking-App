import { useRef, useCallback } from "react";
import { PlaidLink, usePlaidLink } from "react-plaid-link";
import { CSVLink } from "react-csv";

import axios from "axios";

import "../App.css";

// Golang api
const GET_ACCESS_TOKEN_URL =
    "https://ftirs67nxk.execute-api.us-east-1.amazonaws.com/get_access_token";

// // Java api
// const GET_ACCESS_TOKEN_URL = "http://localhost:8080/get_access_token";

// Node.js api
// const GET_ACCESS_TOKEN_URL =
//     "https://d9mv9iihrd.execute-api.us-east-1.amazonaws.com/get_access_token";

// Typescript api
// const GET_ACCESS_TOKEN_URL =
//     "https://12i7pxo3s5.execute-api.us-east-1.amazonaws.com/dev/plaid-node-lambda/get_access_token";

// Local
// const GET_ACCESS_TOKEN_URL = "http://localhost:8080/get_access_token";

// const TRANSACTIONS_GET_URL = "http://localhost:3001/transactions_get";
// const TRANSACTIONS_GET_URL =
//     "https://d9mv9iihrd.execute-api.us-east-1.amazonaws.com/transactions_get";

const XLSX = require("xlsx");
// const path = require('path');
const workSheetName = "가계부";
const filePath = "./가계부.xlsx";

const workSheetColumnName = ["Date", "Usage", "Amount", "Category"];

function Link(props) {
    const { link_token, client_user_id, transactions } = props;

    const exChangePublicTokenForAccessToken = async (public_token) => {
        console.log("public_token: ", public_token);
        axios
            .post(
                GET_ACCESS_TOKEN_URL,
                {
                    public_token: public_token,
                    // uid: client_user_id,
                },
                {
                    headers: {
                        "Content-Type": "application/x-www-form-urlencoded",
                    },
                }
            )
            .then((res) => {
                props.setIsAccessTokenCreated(true);
                console.log(
                    "exChangePublicTokenForAccessToken is Succeed and the access tokein is in your server!!!!"
                );
                console.log("Transaction: ", res.data);
                props.setTransactions(res.data);
            })
            .catch((err) => {
                console.log("get access token failed!!!!!!!!" + err.message);
            });

        // .then((res) => {
        //     axios
        //         .get(TRANSACTIONS_GET_URL)
        //         .then((res) => {
        //             console.log("Getting Transaction succeed!!!!!!!!");
        //             props.setTransactions(res.data);
        //         })
        //         .catch((err) => {
        //             console.log(
        //                 "Getting Transaction failed!!!!!!!!" + err.message
        //             );
        //         });
        // });
    };

    const onSuccess = useCallback((public_token) => {
        console.log("OnSuccess Public token: ", public_token);
        props.setPublicToken(public_token);
        exChangePublicTokenForAccessToken(public_token);
    }, []);

    const config = {
        token: link_token,
        onSuccess,
    };

    const { open, ready } = usePlaidLink(config);

    const openHandler = () => {
        open();
    };

    const generateExcelFile = () => {
        const data = transactions.map((transaction) => {
            return [
                transaction.date,
                transaction.name,
                "$" + transaction.amount,
                transaction.category,
            ];
        });
        const workBook = XLSX.utils.book_new(); // Create a new workbook
        const workSheetData = [workSheetColumnName, ...data];
        const workSheet = XLSX.utils.aoa_to_sheet(workSheetData);
        XLSX.utils.book_append_sheet(workBook, workSheet, workSheetName);
        XLSX.writeFile(workBook, filePath);
        return true;
    };

    return (
        <div>
            <div>
                <button type="button" onClick={() => props.createLinkToken()}>
                    Click this button first (Create Link Token)
                </button>
            </div>

            <div>
                <button type="button" onClick={() => openHandler()}>
                    Get Transactions information
                </button>
            </div>

            <div className="col-sm-8">
                {/* <CSVLink
                    data={transactions}
                    filename="가계부"
                    className="btn btn-success mb-3"
                    onClick={() => {
                        console.log("Export link is clicked!");
                    }}
                >
                    Export CSV File
                </CSVLink>

                <button
                    onClick={() => {
                        generateExcelFile();
                    }}
                >
                    Export Excel File
                </button> */}

                <table className="table table-bordered text-white">
                    <thead>
                        <tr>
                            <th scope="col">Sr. No.</th>
                            <th scope="col">Date</th>
                            <th scope="col">Usage</th>
                            <th scope="col">Amount</th>
                            <th scope="col">Category</th>
                        </tr>
                    </thead>
                    <tbody>
                        {transactions.map((transaction, index) => (
                            <tr key={index}>
                                <td> {index + 1} </td>
                                <td>{transaction.date} </td>
                                <td>{transaction.name} </td>
                                <td>{"$" + transaction.amount} </td>
                                <td>{transaction.category} </td>
                            </tr>
                        ))}
                    </tbody>
                </table>
            </div>
        </div>
    );
}

export default Link;
