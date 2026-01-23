import { ApolloClient, InMemoryCache, HttpLink} from '@apollo/client';

export const client = new ApolloClient({
  link: new HttpLink({ uri: import.meta.env.VITE_API_URL || 'http://localhost:8080/graphql' }),
  cache: new InMemoryCache(),
});