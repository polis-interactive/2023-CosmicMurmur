
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


void main(void)
{

    float t = time;

    vec2 uv = floor(gl_FragCoord.xy / pixel) * pixel / resolution.xy;

    vec3 color = vec3(0.0);

    color += vec3(max(0.0, sin(uv.y - time * 0.1 * PI) * 8.0 - 7.0), 0.0, 0.0);

    color += vec3(0.0, 0.0, max(0.0, sin(uv.y - time * 0.1 * PI + 0.5 * PI) * 10.0 - 9.0));

    color += vec3(0.0, max(0.0, sin(uv.y - time * 0.1 * PI + 1.0 * PI) * 10.0 - 9.0), 0.0);

    color += vec3(max(0.0, sin(uv.y - time * 0.1 * PI + 1.5 * PI) * 10.0 - 9.0));


    gl_FragColor = vec4(color, 1.0);
}