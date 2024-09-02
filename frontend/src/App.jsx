import { useEffect, useState } from 'react'
import reactLogo from './assets/react.svg'
import viteLogo from '/vite.svg'
import './App.css'
import axios from 'axios'

function App() {
  const [count, setCount] = useState(0)
  // useEffect(()=>{
  //   axios.get('/api/metrics').then((resp)=>{
  //     console.log(resp.data)
  //   })
  // },[])

  setInterval(
    function(){
      axios.get('/api/metrics').then((resp)=>{
        console.log(resp.data)
      })
    },1000)
  return (
    <>
      {/* {
        axios.get('/api/metrics').then((resp)=>{
          console.log(resp.data)
        })
      } */}
    </>
  )
}

export default App
