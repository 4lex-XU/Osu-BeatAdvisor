import React, { useState } from "react";

export default function CheckboxLanguages(props) {
  const list = [
    { id: "19", name: "english" },
    { id: "20", name: "chinese" },
    { id: "21", name: "french" },
    { id: "22", name: "german" },
    { id: "23", name: "italian" },
    { id: "24", name: "japanese" },
    { id: "25", name: "korean" },
    { id: "26", name: "spanish" },
    { id: "27", name: "russian" },
    { id: "28", name: "instrumental" },
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
        id="touslangues"
        name="touslangues"
        checked={isCheckAll}
        onChange={handleSelectAll}
      />
      <label htmlFor="touslangues">tous</label>
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
