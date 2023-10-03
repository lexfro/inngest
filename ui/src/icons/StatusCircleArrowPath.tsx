export function IconStatusCircleArrowPath({
  className,
  title,
}: {
  className?: string;
  title?: string;
}) {
  return (
    <svg
      width="28"
      height="28"
      viewBox="0 0 30 30"
      fill="none"
      xmlns="http://www.w3.org/2000/svg"
      className={className}
    >
      <title>{title}</title>
      <g>
        <path
          d="M15 26C21.0751 26 26 21.0751 26 15C26 8.92487 21.0751 4 15 4C8.92487 4 4 8.92487 4 15C4 21.0751 8.92487 26 15 26Z"
          fill="#38BDF8"
        />
        <path
          d="M18.5837 16.0231L18.5837 16.0232C18.2391 17.3871 17.2357 18.4489 15.9557 18.8127C14.6763 19.1764 13.3099 18.7893 12.3707 17.7935L12.3706 17.7934L12.3138 17.7333H13.7693C14.1173 17.7333 14.3896 17.4358 14.3896 17.0814C14.3896 16.7271 14.1173 16.4296 13.7693 16.4296H10.826C10.478 16.4296 10.2057 16.7271 10.2057 17.0814V20.2025C10.2057 20.5569 10.478 20.8544 10.826 20.8544C11.174 20.8544 11.4463 20.5569 11.4463 20.2025V18.6665L11.4885 18.7113C11.4885 18.7113 11.4885 18.7113 11.4885 18.7113C12.74 20.039 14.5664 20.5592 16.2793 20.0724C17.9917 19.5858 19.3256 18.1684 19.7827 16.3585C19.8696 16.0162 19.6798 15.6565 19.3459 15.5612C19.0092 15.4651 18.671 15.6792 18.5837 16.0231ZM19.6065 13.3772L19.6067 13.377C19.7221 13.2543 19.7858 13.0893 19.7858 12.9185V9.79743C19.7858 9.44306 19.5135 9.14561 19.1655 9.14561C18.8175 9.14561 18.5452 9.44306 18.5452 9.79743V11.3342L18.503 11.2894C17.2515 9.96173 15.4251 9.44151 13.7122 9.92831C12.0002 10.4148 10.6665 11.8316 10.2091 13.6409C10.146 13.8656 10.2022 14.11 10.3587 14.2796C10.5163 14.4502 10.7512 14.5182 10.9715 14.4533C11.1908 14.3887 11.356 14.2043 11.4075 13.9762C11.7519 12.6122 12.7553 11.5502 14.0354 11.1863C15.3152 10.8224 16.6822 11.2099 17.6214 12.2065L17.6216 12.2066L17.6777 12.266H16.2229C15.8749 12.266 15.6026 12.5634 15.6026 12.9178C15.6026 13.2721 15.8749 13.5696 16.2229 13.5696H19.1662C19.1662 13.5696 19.1662 13.5696 19.1662 13.5696C19.3325 13.5697 19.4907 13.4997 19.6065 13.3772Z"
          fill="white"
          stroke="white"
          strokeWidth="0.2"
        />
      </g>
    </svg>
  );
}
