import { ApolloClient, InMemoryCache, HttpLink } from "@apollo/client";
import { apiUrl } from "./api";

export const client = new ApolloClient({
  link: new HttpLink({
    uri: apiUrl("/graphql"),
    credentials: "include",
  }),
  cache: new InMemoryCache(),
});
