import React from "react";
import "./style.css";

const ContentWrapper = ({ align, form, children }) => {
  const className =
    align === "top" ? "gm-content-outer top" : "gm-content-outer";
  const maxWidth = form ? 768 : 1200;

  return (
    <div className={className}>
      <div className="gm-content-inner" style={{ maxWidth }}>
        {children}
      </div>
    </div>
  );
};

export default ContentWrapper;
