import React, { useEffect, useState } from 'react'

import { ButtonColors, ButtonSizes } from './Enums'

interface Props {
  size: ButtonSizes
  text: string
  color: ButtonColors
}

export const Button = (props: Props): JSX.Element => {
  const [style, setStyle] = useState<React.CSSProperties>({})

  useEffect(() => {
    switch (props.size) {
      case ButtonSizes.small:
        setStyle({ padding: '0', backgroundColor: props.color })
        break
      case ButtonSizes.medium:
        setStyle({
          padding: '1.5rem 2rem 1.5rem 2rem',
          backgroundColor: props.color,
        })
        break
    }
  }, [])

  return (
    <div>
      <button style={style}>{props.text}</button>
    </div>
  )
}
