import React, { useState } from "react";

export default function CheckboxCategories(props) {
  const list = [
    { id: "01", name: "graveyard" },
    { id: "02", name: "wip" },
    { id: "03", name: "pending" },
    { id: "04", name: "ranked" },
    { id: "05", name: "approved" },
    { id: "06", name: "qualified" },
    { id: "07", name: "loved" },
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
    <div className="row-osu-form">
      <input
        type="checkbox"
        id="touscat"
        name="touscat"
        checked={isCheckAll}
        onChange={handleSelectAll}
      />
      <label htmlFor="touscat">tous</label>
      {list.map(({ id, name }) => (
        <React.Fragment key={id}>
          <input
            id={id}
            type="checkbox"
            name={name}
            onChange={handleClick}
            checked={props.isCheck.includes(name)}
          />
          <label htmlFor={id}>{name}</label>
        </React.Fragment>
      ))}
    </div>
  );
}
