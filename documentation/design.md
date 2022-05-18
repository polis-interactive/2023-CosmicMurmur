
# Cosmic Murmur - Design

## High Level Goals

- Display graphics with high availability
- Accessible for remote monitoring
- Support on the fly configuration changes
- Send alert when unexpected event occurs

## Configuration

- __Lighting__: array for creating segments (array of universe broken tuples, led_count / number_of_strips), number of segments,
    scale
- __Graphics__: global brightness, gamma, shader name
- __Render__: Controller count, parameters

## Types

### Controller

- Universes: array of numbers, all unique
- IP Address: string, unique

### Render

- Shader name: string
- Brightness: number
- Gamma: number

### Lighting

- Width, Height: number
- Construction: Array of universes; universe is an array of tuples; 
    tuples are pairs of led_count / string_count
- Segments: Number
