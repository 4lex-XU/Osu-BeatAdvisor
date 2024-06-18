import RangeSlider from 'react-bootstrap-range-slider';

export default function RangeSliderDiffilculty(props) {
  return (
    <div className="row-osu-form">
      {' '}
      <div className="row-osu-form">
        <label>Difficulté min:</label>

        <RangeSlider
          value={props.min}
          onChange={(e) => {
            if (e.target.value > props.max) {
              props.setMax(e.target.value);
            }
            props.setMin(e.target.value)}}
          min={0}
          max={10}
        />
      </div>
      <div className="row-osu-form">
        <label>Difficulté max:</label>

        <RangeSlider
          value={props.max}
          onChange={(e) => {
            if (e.target.value < props.min) {
              props.setMin(e.target.value);
            }
            props.setMax(e.target.value)}}
          min={0}
          max={10}
        />
      </div>
    </div>
  );
}
