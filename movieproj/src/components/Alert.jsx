export default function Alert(props){
    return(
        <>
        <div className={"alert"+props.className} role="alert">
            {props.message }
        </div>

        </>
    )
}