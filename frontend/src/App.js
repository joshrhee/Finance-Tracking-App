import "./App.css";

import Link from "./Components/Link";
import { useCallback, useEffect, useState } from "react";

import axios from "axios";
import { v4 as uuidv4 } from "uuid";

// Golang api
const CREATE_LINK_TOKEN_URL =
    "https://ftirs67nxk.execute-api.us-east-1.amazonaws.com/create_link_token";

// Node.js api
// const CREATE_LINK_TOKEN_URL =
//     "https://d9mv9iihrd.execute-api.us-east-1.amazonaws.com/create_link_token";

// Typescript api
// const CREATE_LINK_TOKEN_URL =
//     "https://12i7pxo3s5.execute-api.us-east-1.amazonaws.com/dev/plaid-node-lambda/create_link_token";

// Local
// const CREATE_LINK_TOKEN_URL = "http://localhost:8080/create_link_token";

function App() {
    const [link_token, setLink_token] = useState("");
    // const [link_token_expiration, setLink_token_expiration] = useState("");
    const [public_token, setPublic_token] = useState("");

    const [isAccessTokenCreated, setIsAccessTokenCreated] = useState(false);
    const [transactions, setTransactions] = useState([]);

    const CLIENT_USER_ID = uuidv4(); // 이 uuid는 나만씀, 결국 다른 사람들도 이거 쓰려면 로그인 정보가지고 uuid만들어야함!!!

    const createLinkToken = async () => {
        axios
            .post(CREATE_LINK_TOKEN_URL, { CLIENT_USER_ID })
            .then((res) => {
                console.log("Successfully created link token!!!");
                setLink_token(res.data.link_token);
                // setLink_token_expiration(res.data.expiration);

                return res.data.link_token;
            })
            .catch((err) => {
                console.log("Creating link token failed!!!!");
                console.log("Error: ", err);
            });
    };

    const setTransactionHandler = useCallback((transactions) => {
        setTransactions(transactions);
    }, []);

    const setPublicTokenHandler = useCallback((public_token) => {
        setPublic_token(public_token);
    }, []);

    useEffect(() => {
        console.log("linktoken: ", link_token);
        // console.log("link token expiration: ", link_token_expiration);
    }, [link_token]);

    useEffect(() => {
        console.log("public_token: ", public_token);
    }, [public_token]);

    useEffect(() => {
        console.log("isAccessTokenCreated is chagend!!");
    }, [isAccessTokenCreated]);

    useEffect(() => {
        console.log("Transactions: ", transactions);
    }, [transactions]);

    return (
        <div className="App">
            <div>
                <button onClick={createLinkToken}>
                    Click this button first (Create Link Token)
                </button>
            </div>
            <div className="Link">
                <Link
                    createLinkToken={createLinkToken}
                    setPublicToken={setPublicTokenHandler}
                    setIsAccessTokenCreated={setIsAccessTokenCreated}
                    link_token={link_token}
                    client_user_id={CLIENT_USER_ID}
                    transactions={transactions}
                    setTransactions={setTransactionHandler}
                />
            </div>
        </div>
    );
}

export default App;
