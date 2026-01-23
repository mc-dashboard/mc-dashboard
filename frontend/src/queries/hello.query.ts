import { gql } from "@apollo/client"
import type { TypedDocumentNode } from "@apollo/client"
import type { GetTodoQuery, GetTodoQueryVariables } from "../types/__generated__/graphql"

export const GET_TODO: TypedDocumentNode<
GetTodoQuery,
GetTodoQueryVariables
> = gql`
    query getTodo {
        todo {
            text
            done
            user {
            name
            }
        }
    }
`