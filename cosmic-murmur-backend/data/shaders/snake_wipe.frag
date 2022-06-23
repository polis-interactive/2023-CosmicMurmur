
#ifdef GL_ES

precision highp float;
#define IN varying
#define OUT out
#define TEXTURE texture2D

#else

#define IN in
#define OUT out
#define TEXTURE texture

#endif

uniform float time;
uniform float pixel;
uniform vec2 resolution;

#define PI 3.14159265358979323846
vec3 c[4] = vec3[](vec3(1.0, 0.0, 0.0),vec3(0.0, 0.0, 1.0), vec3(0.0, 1.0, 0.0), vec3(1.0, 1.0, 1.0));

float parabola( float x, float k ){
    return pow( 4.0*x*(1.0-x), k );
}

vec3 hsb2rgb( in vec3 c ){
    vec3 rgb = clamp(abs(mod(c.x*6.0+vec3(0.0,4.0,2.0),
    6.0)-3.0)-1.0,
    0.0,
    1.0 );
    rgb = rgb*rgb*(3.0-2.0*rgb);
    return c.z * mix(vec3(1.0), rgb, c.y);
}



void main(void)
{

	vec3 color = vec3(0.0);

    vec2 grid = resolution.xy / pixel;

    vec2 uv_grid = floor(gl_FragCoord.xy / pixel);

    float snake_number = grid.y * uv_grid.x + (mod(uv_grid.x, 2.0) == 0.0 ? uv_grid.y : grid.y - uv_grid.y - 1.0);

    for ( float i = 0.0; i < 4.0; ++i )
    {
        float t = mod(time + grid.x / 4.0 * i, grid.x);
        float y = max(0.0, parabola(snake_number / 12.0 - t, 1.0));
        color += c[uint(i)] * y;
    }



    gl_FragColor = vec4(color, 1.0);
}