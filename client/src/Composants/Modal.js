import React from "react";
import "../CSS/styles.css";

function Modal({ isOpen, onClose, children }) {
  if (!isOpen) return null;

  return (
    <div className="modal">
      <div className="modal-content">
        {children}
        <button onClick={onClose}>Fermer</button>
      </div>
    </div>
  );
}

export default Modal;
