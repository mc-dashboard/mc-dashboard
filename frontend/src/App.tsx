import { useState } from 'react'
import { useQuery } from '@apollo/client/react'
import { GET_TODO } from './queries/hello.query'

function App() {
  const [count, setCount] = useState(0)
  const {loading, error, data} = useQuery(GET_TODO)
  if (error || !data) return <p>Error : {error?.message}</p>;
  console.log(data);

  return (
    <>
    <p>message from graphql: {data?.todo.text}</p>

      <div className="card">
        <button onClick={() => setCount((count) => count + 1)}>
          count is {count}
        </button>
      </div>
    </>
  )
}

export default App
