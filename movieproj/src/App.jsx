import Homepage from './components/Homepage'
import './App.css'
import { Outlet, useNavigate } from 'react-router-dom'
import { Link } from 'react-router-dom'
import Alert from './components/Alert'
import { useEffect, useState } from 'react'
const navigate=useNavigate
function App() {
  const [jwttok, setjwttok] = useState("he")
  const [alertmsg,setalertmsg]=useState('')
  const [alertclass,setalertclass]=useState('d-none')
const logout=()=>{
  setjwttok('')
  navigate("/login")

}
useEffect(() => {
  const refreshToken = async () => {
    if (jwttok === "") {
      try {
        const response = await fetch("http://localhost:8080/refresh", {
          method: "GET",
          credentials: "include"
        });
        
        if (!response.ok) {
          throw new Error('Refresh failed');
        }
        
        const data = await response.json();
        if (data.access_token) {
          setjwttok(data.access_token);
        }
      } catch (error) {
        console.error("Error refreshing token:", error);
        
      }
    }
  };

  refreshToken();
}, [jwttok]);
  return (
    <>
      <div className="hdrpbtn">
        <h1 className="header">Movie with Go</h1>
        {!jwttok ? <Link to="/login"><button className='Signbtn'>Sign in</button></Link> : <Link to="/#"><button className='badge bg-danger' onClick={logout}>Logout</button></Link>}

      </div>
      <hr />
      <div className='sidebar'>
        <Link to="/"> <div className=" innav  Home"> Home</div></Link>
        <Link to="/movies"><div className=" innav  Movies">Movies</div></Link>
        <Link to="/genre"><div className=" innav  Games">Genre</div></Link>
        {jwttok &&
        <>
        <Link to="/admin/movies/0"><div className=" innav  Add Movie">Add Movie</div></Link>
        <Link to="/manage"><div className=" innav  Manage">Manage</div></Link>
        <Link to="/graphql"><div className=" innav  GraphQl">Graphql</div></Link>
        </>
}
      </div>
      <Alert message={alertmsg}
      className={alertclass}/>
      <Outlet  context={{jwttok,setjwttok,setalertclass,setalertmsg}}/>
    </>
  )
}

export default App
