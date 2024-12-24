import './homepage.css'
import Carousel from './Carousel'
import tickimg from './movie_tickets.jpg'
export default  function Homepage(){
    return(
        <>
        <div className="main">
            <Carousel />
        </div>
        <hr />
        <div className='randcont'>
            <h1 className='inrandcont'>Find a Movie to Watch Tonight</h1>
            <div className='movieticksimg'>
            <img src={tickimg} alt="tickimg" />
            <img src={tickimg} alt="tickimg" />
            <img src={tickimg} alt="tickimg" />
            <img src={tickimg} alt="tickimg" />
            </div>
        </div>
        </>
    )
}