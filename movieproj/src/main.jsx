import { StrictMode } from 'react'
import { createRoot } from 'react-dom/client'
import App from './App.jsx'
import './index.css'
import { createBrowserRouter, RouterProvider } from 'react-router-dom'
import ErrorPage from './components/ErrorPage.jsx'
import Homepage from './components/Homepage.jsx'
import Movies from './components/Movies.jsx'
import Genre from './components/Genre.jsx'
import Login from './components/Login.jsx'
import EditMovie from './components/EditMovie.jsx'
import Manage from './components/Manage.jsx'
import Graphql from './components/Graphql.jsx'
import Movie from './components/Movie.jsx'
const router = createBrowserRouter([
  {
    path: '/',
    element: <App />,
    errorElement: <ErrorPage />,
    children: [
      { index: true, element: <Homepage /> },
      {
        path: '/movies',
        element: <Movies />,
      },
      {
        path:'/movies/:id',
        element:<Movie/>
      },
      {
        path: '/genre',
        element: <Genre />,
      },
      {
        path:'/login',
        element:<Login/>
      },
      {
        path:'/admin/movies/0',
        element:<EditMovie/>
      },
      {
        path:'/admin/movie/:id',
        element:<EditMovie/>
      },
      {
        path:'/manage',
        element:<Manage/>
      },
      {
        path:'/graphql',
        element:<Graphql/>
      }
    ],
  },
]);

createRoot(document.getElementById('root')).render(
  <StrictMode>
 <RouterProvider router={router}/>
  </StrictMode>,
)
