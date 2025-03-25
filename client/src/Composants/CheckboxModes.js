import React, { useState } from "react";

export default function CheckboxModes(props) {
  const list = [
    { id: "29", name: "osu" },
    { id: "30", name: "taiko" },
    { id: "31", name: "fruits" },
    { id: "32", name: "mania" },
  ];
  const [isCheckAll, setIsCheckAll] = useState(false);

  const handleSelectAll = (e) => {
    setIsCheckAll(!isCheckAll);
    props.setIsCheck(list.map((li) => li.name));
    if (isCheckAll) {
      props.setIsCheck([]);
    }
  };

  const handleClick = (e) => {
    const { name, checked } = e.target;
    props.setIsCheck([...props.isCheck, name]);
    if (!checked) {
      props.setIsCheck(props.isCheck.filter((item) => item !== name));
    }
  };
  return (
      <div className="checkbox-group modes">
        <label className="checkbox-item modes-item">
          <input
              type="checkbox"
              id="tousmode"
              name="tousmode"
              checked={isCheckAll}
              onChange={handleSelectAll}
          />
          tous
        </label>
        {list.map(({ id, name }) => (
            <label className="checkbox-item modes-item" key={id}>
              <input
                  id={id}
                  type="checkbox"
                  name={name}
                  onChange={handleClick}
                  checked={props.isCheck.includes(name)}
              />
              {name}
            </label>
        ))}
      </div>
  );

}
