import React, { useState } from "react";
import { Button, InlineFormLabel } from "@grafana/ui";

type SecretMessageProps = {
  message: string;
};

export const SecretMessage = (props: SecretMessageProps) => {
  const [show, setShow] = useState(false);
  const { message = "hello" } = props;
  const secret = message
    .split("")
    .map(() => "*")
    .join("");
  return (
    <>
      <div className="gf-form">
        <InlineFormLabel width={12}>Secret Message</InlineFormLabel>
        <div style={{ padding: "8px" }}>{show ? message : secret}</div>
        <Button
          variant="secondary"
          size="sm"
          icon={show ? "eye-slash" : "eye"}
          style={{ margin: "8px" }}
          onClick={(e) => {
            setShow(!show);
            e.preventDefault();
          }}
        >
          {show ? "hide" : "reveal"}
        </Button>
      </div>
    </>
  );
};
