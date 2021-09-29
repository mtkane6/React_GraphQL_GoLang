const Input = (props) =>{
    return (
        <div className="mb-3">
            <label htmlFor={props.name} className="form-label">
                {props.title}
            </label>
            <input
                type={props.type}
                className={`form-control ${props.className}`} // form-control will always exist, the other one only if passed in props.
                id={props.name}
                name={props.name}
                value={props.value}
                onChange={props.handleChange}
                placeholder={props.placeholder}
            />
            <div className={props.errorDiv}> {/* this area for form error handling */}
                {props.errorMsg}
            </div>
        </div>
    );
};

export default Input;