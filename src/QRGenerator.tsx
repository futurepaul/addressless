import { useState, useLayoutEffect } from "react";
import QRCode from "qrcode";

function QRGenerator({ qrCode, width }: { qrCode: string, width?: number}) {
  const [code, setCode] = useState(qrCode);

  useLayoutEffect(() => {
    QRCode.toDataURL(qrCode, {
      errorCorrectionLevel: "L",
      width: width
    }).then(setCode);
  }, [qrCode]);

  return <img src={code} alt="QR Code Invoice" />;
}

export default QRGenerator;