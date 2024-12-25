import { useState } from "react"
import './Logiform.css'
import { useNavigate, useOutletContext } from "react-router-dom"
export default function Login() {
    const [email, setemail] = useState("")
    const [password, setPassword] = useState("")
  const {setjwttok}=useOutletContext()
  const {setalertmsg}=useOutletContext()
  const {setalertclass}=useOutletContext()
  const {toggleRefresh}=useOutletContext()
  const navigate=useNavigate()
    const handlesubmit=(e)=>{
        e.preventDefault()
       let payload={
        email:email,
        password:password
       }

       const requestoptions={
        method: "POST",
        headers: { 'Content-Type': 'application/json' },
        credentials:"include",
        body: JSON.stringify(payload),
       }
       
       fetch("http://localhost:8080/authenticate",requestoptions)
       .then(response=>response.json())
       .then(data=>{
        if(data.error){
            setalertclass("alert-danger")
            setalertmsg("Invalid Credentials")
            console.log("error in auth block",data.error)
        }else{
            setjwttok(data.access_token)
            setalertclass("alert-success")
            toggleRefresh(true)
            navigate("/")
        }
       })
       .catch(error=>{
        setalertclass("alert-danger")
        setalertmsg("Error Occured")
        console.log(error)
       })

    }
    return (
        <>
            <h1>Login Page</h1>
            <hr />
            <div className="Form">
                <form action="" onSubmit={handlesubmit}>
                    <label htmlFor="username">Email:</label>
                    <input type="email" value={email} onChange={(e) => {
                        setemail(e.target.value)

                    }} />
                    <br />
                    <label htmlFor="password">Password:</label>
                    <input type="password" value={password} onChange={(e)=>{
                        setPassword(e.target.value)
                    }}/>
                    <input  type="submit" value="Login" />




                </form>
            </div>
        </>

    )
}