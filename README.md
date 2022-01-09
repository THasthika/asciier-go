# asciier - ASCII Art Tool

A tool to convert images to and from ascii art

## ASCII Mapping Used

```go
switch {
  case pixelValue < 0.1:
		return '.'
	case pixelValue < 0.2:
		return ','
	case pixelValue < 0.3:
		return ';'
	case pixelValue < 0.4:
		return '!'
	case pixelValue < 0.5:
		return 'v'
	case pixelValue < 0.6:
		return 'l'
	case pixelValue < 0.7:
		return 'L'
	case pixelValue < 0.8:
		return 'F'
	case pixelValue < 0.9:
		return 'E'
	default:
		return '$'
}
```

## Tasks

- [x] Convert Image to ASCII txt file
- [ ] Convert ASCII txt file to grayscale image

## Usage

```bash
asciier [options] [image_file]

options:
	-w [width]
	-h [height]
	-o [output_file]
```

## Example Usage

![](./example/ex1.gif)
