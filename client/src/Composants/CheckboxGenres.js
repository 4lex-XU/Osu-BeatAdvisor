import React, { useState } from 'react';

export default function CheckboxGenres(props) {
  const list = [
    { id: '08', name: 'video game' },
    { id: '09', name: 'anime' },
    { id: '10', name: 'rock' },
    { id: '11', name: 'pop' },
    { id: '12', name: 'funk' },
    { id: '13', name: 'hip hop' },
    { id: '14', name: 'electronic' },
    { id: '15', name: 'metal' },
    { id: '16', name: 'classical' },
    { id: '17', name: 'folk' },
    { id: '18', name: 'jazz' },
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
        id="tousgenre"
        name="tousgenre"
        checked={isCheckAll}
        onChange={handleSelectAll}
      />
      <label htmlFor="tousgenre">tous</label>
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
