const Select = (props) => {
    return (
        <div className="mb-3">
            <label htmlFor={props.name} className="form-label">
                {""+props.title}
            </label>
            <select 
                name={props.name}  /* <- this matches with handleChange().  '[name]'. Must match spelling in this.state.*/ 
                className="form-select" 
                value={props.value} 
                onChange={props.handleChange}>

                <option value="">{props.placeholder}</option>
                {props.options.map(
                    (option) => {
                        return (
                            <option className="form-select"
                                key={option.id}
                            >
                                {option.value}
                            </option>
                        )
                    }
                )}
                {/* <option className="form-select" value={"G"}>G</option>
                <option className="form-select" value={"PG"}>PG</option>
                <option className="form-select" value={"PG-13"}>PG-13</option>
                <option className="form-select" value={"R"}>R</option>
                <option className="form-select" value={"NC-17"}>NC-17</option> */}
            </select>
        </div>
    );
}

export default Select;