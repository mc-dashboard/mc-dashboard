import { ApolloClient, InMemoryCache, HttpLink } from "@apollo/client";
import { API_BASE_URL } from "./api";

export const client = new ApolloClient({
  link: new HttpLink({
    uri: `${API_BASE_URL}/graphql`,
    credentials: "include",
  }),
  cache: new InMemoryCache(),
});
