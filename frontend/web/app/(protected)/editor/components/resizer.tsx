type Props = {
  initDrag: (e: React.MouseEvent) => void;
};

export default function Resizer({ initDrag }: Props) {
  return <div className="resizer" onMouseDown={initDrag} />;
}
