package store

import "time"

const (
	importData = `{ "_id" : 10140, "created_at" : { "$date" : 1374639571037 }, "image_url" : "/thumbs/10140.png", "modified_at" : { "$date" : 1374639571037 }, "parent" : 9970, "parent_version" : 0, "user" : "ca2d2b8", "versions" : [ { "created_at" : { "$date" : 1374639571037 }, "code" : "#ifdef GL_ES\nprecision mediump float;\n#endif\n\nuniform float time;\nuniform vec2 mouse;\nuniform vec2 resolution;\nfloat iGlobalTime = time;\n\n// Mountains. (C) David Hoskins - 2013\n\n\n// https://www.shadertoy.com/view/4slGD4\n\n// A ray-marched version of my terrain renderer which uses\n// streaming texture normals for speed:-\n// http://www.youtube.com/watch?v=qzkBnCBpQAM\n\n// It was difficult finding a suitable starting point, but I think this one works OK\n\n// It uses binary subdivision to accurately find the height map.\n// Lots of thanks to Iñigo and his noise functions!\n\n// Video of my OpenGL version that \n// http://www.youtube.com/watch?v=qzkBnCBpQAM\n\n// Stereo version code thanks to Croqueteer :)\n//#define STEREO \n\n// Remove the following line to take out the trees.\n#define TREES\n\n#ifdef TREES\nfloat treeLine = 0.0;\nfloat treeCol =100.0;\n#endif\n\nvec3 sunLight  = normalize( vec3(  0.4, 0.4,  0.48 ) );\nvec3 sunColour = vec3(1.0, .9, .83);\nfloat specular = 0.0;\nvec3 cameraPos;\n\n\n// This peturbs the fractal positions for each iteration down...\n// Helps make nice twisted landscapes...\nconst mat2 rotate2D = mat2(1.4623, 1.67231, -1.67231, 1.4623);\n\n// Alternative rotation:-\n// const mat2 rotate2D = mat2(1.2323, 1.999231, -1.999231, 1.22);\n\n//--------------------------------------------------------------------------\n// Noise functions...\nfloat Hash( float n )\n{\n    return fract(sin(n)*43758.5453123);\n}\n\n//--------------------------------------------------------------------------\nfloat Hash(vec2 p)\n{\n\treturn fract(sin(dot(p, vec2(12.9898, 78.233))) * 43758.5453);\n}\n\n//--------------------------------------------------------------------------\nfloat Noise( in vec3 x )\n{\n    vec3 p = floor(x);\n    vec3 f = fract(x);\n    f = f*f*(3.0-2.0*f);\n    float n = p.x + p.y*57.0 + 113.0*p.z;\n    float res = mix(mix(mix( Hash(n+  0.0), Hash(n+  1.0),f.x),\n                        mix( Hash(n+ 57.0), Hash(n+ 58.0),f.x),f.y),\n                    mix(mix( Hash(n+113.0), Hash(n+114.0),f.x),\n                        mix( Hash(n+170.0), Hash(n+171.0),f.x),f.y),f.z);\n    return res;\n}\n//--------------------------------------------------------------------------\nfloat Noise( in vec2 x )\n{\n    vec2 p = floor(x);\n    vec2 f = fract(x);\n    f = f*f*(3.0-2.0*f);\n    float n = p.x + p.y*57.0;\n    float res = mix(mix( Hash(n+  0.0), Hash(n+  1.0),f.x),\n                    mix( Hash(n+ 57.0), Hash(n+ 58.0),f.x),f.y);\n    return res;\n}\n\n//--------------------------------------------------------------------------\nvec2 Noise2( in vec2 x )\n{\n\tvec2 res = vec2(Noise(x), Noise(x+vec2(4101.03, 2310.0)));\n    return res-vec2(.5, .5);\n}\n\n//--------------------------------------------------------------------------\n// iq's derivative noise function...\nvec3 NoiseDerivative( in vec2 x )\n{\n    vec2 p = floor(x);\n    vec2 f = fract(x);\n    vec2 u = f*f*(3.0-2.0*f);\n    float n = p.x + p.y*57.0;\n    float a = Hash(n+  0.0);\n    float b = Hash(n+  1.0);\n    float c = Hash(n+ 57.0);\n    float d = Hash(n+ 58.0);\n\treturn vec3(a+(b-a)*u.x+(c-a)*u.y+(a-b-c+d)*u.x*u.y,\n\t\t\t\t30.0*f*f*(f*(f-2.0)+1.0)*(vec2(b-a,c-a)+(a-b-c+d)*u.yx));\n}\n\n//--------------------------------------------------------------------------\n#ifdef TREES\nfloat Trees(vec2 p)\n{\n\tp *= 5.0;\n\tvec2 rnd = Noise2(p);//vec2(Hash(floor(p.x*4.0)), Hash(floor(p.y*4.0)))*.5;\n\tvec2 v2 = fract(p+rnd)-.5;\n\treturn max(.5-(length(v2)), 0.0) * treeLine*.6;\n}\n#endif\n\n//--------------------------------------------------------------------------\n// Low def version for ray-marching through the height field...\nfloat Terrain( in vec2 p)\n{\n\tvec2 pos = p*0.08;\n\tfloat w = (Noise(pos*.25)*0.75+.15);\n\tw = 36.0 * w * w;\n\tvec2 dxy = vec2(0.0, 0.0);\n\tfloat f = .0;\n\tfor (int i = 0; i < 5; i++)\n\t{\n\t\tvec3 v = NoiseDerivative(pos);\n\t\tdxy += v.yz;\n\t\tf += (w * (v.x) / (1.0 + dot(dxy, dxy))) ;\n\t\tw = -w * 0.37;\t//...Flip negative and positive for variation\n\t\tpos = rotate2D * pos;\n\t}\n\tfloat ff = Noise(pos*.003);\n\t\n\tf += pow(ff, 6.0)*85.-1.0;\n\treturn f;\n}\n\n//--------------------------------------------------------------------------\n// Map to lower resolution for height field mapping for Scene function...\nfloat Map(in vec3 p)\n{\n\tfloat h = Terrain(p.xz);\n\t\t\n\t#ifdef TREES\n\tfloat ff = Noise(p.xz*1.3)*.8;\n\ttreeLine = smoothstep(ff, .1+ff, h) * smoothstep(.5+ff, .4+ff, h);\n\ttreeCol = Trees(p.xz);\n\th += treeCol;\n\t#endif\n\t\n    return p.y - h;\n}\n\n//--------------------------------------------------------------------------\n// High def version only used for grabbing normal information.\nfloat Terrain2( in vec2 p)\n{\n\t// There's some real magic numbers in here! \n\t// The Noise calls add large mountain ranges for more variation over distances...\n\tvec2 pos = p*0.08;\n\tfloat w = (Noise(pos*.25)*0.75+.15);\n\tw = 36.0 * w * w;\n\tvec2 dxy = vec2(0.0, 0.0);\n\tfloat f = .0;\n\tfor (int i = 0; i < 5; i++)\n\t{\n\t\tvec3 v = NoiseDerivative(pos);\n\t\tdxy += v.yz;\n\t\tf += (w * (v.x)  / (1.0 + dot(dxy, dxy)));\n\t\tw =  - w * 0.37;\t//...Flip negative and positive for varition\t   \n\t\tpos = rotate2D * pos;\n\t}\n\tfloat ff = Noise(pos*.003);\n\tf += pow(ff, 6.0)*85.-1.0;\n\t\n\t#ifdef TREES\n\ttreeCol = Trees(p);\n\tf += treeCol;\n\tif (treeCol > 0.0) return f;\n\t#endif\n\t\n\t// That's the last of the low resolution, now go down further for the Normal data...\n\tfor (int i = 0; i < 6; i++)\n\t{\n\t\tvec3 v = NoiseDerivative(pos);\n\t\tdxy += v.yz;\n\t\tf += (w * (v.x) / (1.0 + dot(dxy, dxy)));\n\t\tw =  - w * 0.37;\n\t\tpos = rotate2D * pos;\n\t}\n\t\n\t\n\treturn f;\n}\n\n//--------------------------------------------------------------------------\nfloat FractalNoise(in vec2 xy)\n{\n\tfloat w = .65;\n\tfloat f = 0.0;\n\n\tfor (int i = 0; i < 4; i++)\n\t{\n\t\tf += Noise(xy) * w;\n\t\tw *= 0.5;\n\t\txy *= 2.3;\n\t}\n\treturn f;\n}\n\n//--------------------------------------------------------------------------\n// Simply Perlin clouds that fade to the horizon...\n// 200 units above the ground...\nvec3 GetClouds(in vec3 sky, in vec3 rd)\n{\n\tif (rd.y < 0.0) return sky;\n\tfloat v = (200.0-cameraPos.y)/rd.y;\n\trd.xz *= v;\n\trd.xz += cameraPos.xz;\n\trd.xz *= .010;\n\tfloat f = (FractalNoise(rd.xz) -.55) * 5.0;\n\t// Uses the ray's y component for horizon fade of fixed colour clouds...\n\tsky = mix(sky, vec3(.55, .55, .52), clamp(f*rd.y-.1, 0.0, 1.0));\n\n\treturn sky;\n}\n\n\n\n//--------------------------------------------------------------------------\n// Grab all sky information for a given ray from camera\nvec3 GetSky(in vec3 rd)\n{\n\tfloat sunAmount = max( dot( rd, sunLight), 0.0 );\n\tfloat v = pow(1.0-max(rd.y,0.0),5.)*.5;\n\tvec3  sky = vec3(v*sunColour.x*0.4+0.18, v*sunColour.y*0.4+0.22, v*sunColour.z*0.4+.4);\n\t// Wide glare effect...\n\tsky = sky + sunColour * pow(sunAmount, 6.5)*.32;\n\t// Actual sun...\n\tsky = sky+ sunColour * min(pow(sunAmount, 1150.0), .3)*.65;\n\treturn sky;\n}\n\n//--------------------------------------------------------------------------\n// Merge mountains into te sky background for correct disappearance...\nvec3 ApplyFog( in vec3  rgb, in float dis, in vec3 dir)\n{\n\tfloat fogAmount = clamp(dis* 0.0000165, 0.0, 1.0);\n\treturn mix( rgb, GetSky(dir), fogAmount );\n}\n\n//--------------------------------------------------------------------------\n// Calculate sun light...\nvoid DoLighting(inout vec3 mat, in vec3 pos, in vec3 normal, in vec3 eyeDir, in float dis)\n{\n\tfloat h = dot(sunLight,normal);\n\tfloat c = max(h, 0.0)+.1;\n\tmat = mat * sunColour * c ;\n\t// Specular...\n\tif (h > 0.0)\n\t{\n\t\tvec3 R = reflect(sunLight, normal);\n\t\tfloat specAmount = pow( max(dot(R, normalize(eyeDir)), 0.0), 3.0)*specular;\n\t\tmat = mix(mat, sunColour, specAmount);\n\t}\n}\n\n//--------------------------------------------------------------------------\n// Hack the height, position, and normal data to create the coloured landscape\nvec3 TerrainColour(vec3 pos, vec3 normal, float dis)\n{\n\tvec3 mat;\n\tspecular = .0;\n\tvec3 dir = normalize(pos-cameraPos);\n\t\n\tvec3 matPos = pos * 2.0;// ... I had change scale halfway though, this lazy multiply allow me to keep the graphic scales I had\n\n\tfloat disSqrd = dis * dis;// Squaring it gives better distance scales.\n\n\tfloat f = clamp(Noise(matPos.xz*.05), 0.0,1.0);//*10.8;\n\tf += Noise(matPos.xz*.1+normal.yz*1.08)*.85;\n\tf *= .55;\n\tvec3 m = mix(vec3(.63*f+.2, .7*f+.1, .7*f+.1), vec3(f*.43+.1, f*.3+.2, f*.35+.1), f*.65);\n\tmat = m*vec3(f*m.x+.36, f*m.y+.30, f*m.z+.28);\n\t// Should have used smoothstep to add colours, but left it using 'if' for sanity...\n\tif (normal.y < .5)\n\t{\n\t\tfloat v = normal.y;\n\t\tfloat c = (.5-normal.y) * 4.0;\n\t\tc = clamp(c*c, 0.1, 1.0);\n\t\tf = Noise(vec2(matPos.x*.09, matPos.z*.095+matPos.yy*0.15));\n\t\tf += Noise(vec2(matPos.x*2.233, matPos.z*2.23))*0.5;\n\t\tmat = mix(mat, vec3(.4*f), c);\n\t\tspecular+=.1;\n\t}\n\n\t// Grass. Use the normal to decide when to plonk grass down...\n\tif (matPos.y < 45.35 && normal.y > .65)\n\t{\n\n\t\tm = vec3(Noise(matPos.xz*.073)*.5+.15, Noise(matPos.xz*.12)*.6+.25, 0.0);\n\t\tm *= (normal.y- 0.75)*.85;\n\t\tmat = mix(mat, m, clamp((normal.y-.65)*1.3 * (45.35-matPos.y)*0.1, 0.0, 1.0));\n\t}\n\t#ifdef TREES\n\tif (treeCol > 0.0)\n\t{\n\t\tmat = vec3(.02+Noise(matPos.xz*5.0)*.03, .05, .0);\n\t\tnormal = normalize(normal+vec3(Noise(matPos.xz*33.0)*1.0-.5, .0, Noise(matPos.xz*33.0)*1.0-.5));\n\t\tspecular = .0;\n\t}\n\t#endif\n\t\n\t// Snow topped mountains...\n\tif (matPos.y > 50.0 && normal.y > .28)\n\t{\n\t\tfloat snow = clamp((matPos.y - 50.0 - Noise(matPos.xz * .1)*28.0) * 0.035, 0.0, 1.0);\n\t\tmat = mix(mat, vec3(.7,.7,.8), snow);\n\t\tspecular += snow;\n\t}\n\t// Beach effect...\n\tif (matPos.y < 1.45)\n\t{\n\t\tif (normal.y > .4)\n\t\t{\n\t\t\tf = Noise(matPos.xz * .084)*1.5;\n\t\t\tf = clamp((1.45-f-matPos.y) * 1.34, 0.0, .67);\n\t\t\tfloat t = (normal.y-.4);\n\t\t\tt = (t*t);\n\t\t\tmat = mix(mat, vec3(.09+t, .07+t, .03+t), f);\n\t\t}\n\t\t// Cheap under water darkening...it's wet after all...\n\t\tif (matPos.y < 0.0)\n\t\t{\n\t\t\tmat *= .5;\n\t\t}\n\t}\n\n\tDoLighting(mat, pos, normal,dir, disSqrd);\n\t\n\t// Do the water...\n\tif (cameraPos.y < 0.0)\n\t{\n\t\t// Can go under water, but current camera doesn't find a place...\n\t\tmat = mix(mat, vec3(0.0, .1, .2), .75); \n\t}else\n\tif (matPos.y < 0.0)\n\t{\n\t\t// Pull back along the ray direction to get water surface point at y = 0.0 ...\n\t\tfloat time = (iGlobalTime)*.03;\n\t\tvec3 watPos = matPos;\n\t\twatPos += -dir * (watPos.y/dir.y);\n\t\t// Make some dodgy waves...\n\t\tfloat tx = cos(watPos.x*.052) *4.5;\n\t\tfloat tz = sin(watPos.z*.072) *4.5;\n\t\tvec2 co = Noise2(vec2(watPos.x*4.7+1.3+tz, watPos.z*4.69+time*35.0-tx));\n\t\tco += Noise2(vec2(watPos.z*8.6+time*13.0-tx, watPos.x*8.712+tz))*.4;\n\t\tvec3 nor = normalize(vec3(co.x, 20.0, co.y));\n\t\tnor = normalize(reflect(dir, nor));//normalize((-2.0*(dot(dir, nor))*nor)+dir);\n\t\t// Mix it in at depth transparancy to give beach cues..\n\t\tmat = mix(mat, GetClouds(GetSky(nor), nor), clamp((watPos.y-matPos.y)*1.1, .4, .66));\n\t\t// Add some extra water glint...\n\t\tfloat sunAmount = max( dot(nor, sunLight), 0.0 );\n\t\tmat = mat + sunColour * pow(sunAmount, 228.5)*.6;\n\t}\n\tmat = ApplyFog(mat, disSqrd, dir);\n\treturn mat;\n}\n\n//--------------------------------------------------------------------------\nfloat BinarySubdivision(in vec3 rO, in vec3 rD, float t, float oldT)\n{\n\t// Home in on the surface by dividing by two and split...\n\tfor (int n = 0; n < 4; n++)\n\t{\n\t\tfloat halfwayT = (oldT + t ) * .5;\n\t\tvec3 p = rO + halfwayT*rD;\n\t\tif (Map(p) < 0.25)\n\t\t{\n\t\t\tt = halfwayT;\n\t\t}else\n\t\t{\n\t\t\toldT = halfwayT;\n\t\t}\n\t}\n\treturn t;\n}\n\n//--------------------------------------------------------------------------\nbool Scene(in vec3 rO, in vec3 rD, out float resT )\n{\n    float t = 1.2;\n\tfloat oldT = 0.0;\n\tfloat delta = 0.0;\n\tfor( int j=0; j<170; j++ )\n\t{\n\t\tif (t > 240.0) return false; // ...Too far\n\t    vec3 p = rO + t*rD;\n        if (p.y > 95.0) return false; // ...Over highest mountain\n\n\t\tfloat h = Map(p); // ...Get this positions height mapping.\n\t\t// Are we inside, and close enough to fudge a hit?...\n\t\tif( h < 0.25)\n\t\t{\n\t\t\t// Yes! So home in on height map...\n\t\t\tresT = BinarySubdivision(rO, rD, t, oldT);\n\t\t\treturn true;\n\t\t}\n\t\t// Delta ray advance - a fudge between the height returned\n\t\t// and the distance already travelled.\n\t\t// It's a really fiddly compromise between speed and accuracy\n\t\t// Too large a step and the tops of ridges get missed.\n\t\tdelta = max(0.01, 0.2*h) + (t*0.0065);\n\t\toldT = t;\n\t\tt += delta;\n\t}\n\n\treturn false;\n}\n\n//--------------------------------------------------------------------------\nvec3 CameraPath( float t )\n{\n\tfloat m = 1.0+(mouse.x)*300.0;\n\tt = (iGlobalTime*1.5+m+350.)*.006 + t;\n    vec2 p = 276.0*vec2( sin(3.5*t), cos(1.5*t) );\n\treturn vec3(140.0-p.x, 0.6, -88.0+p.y);\n}\n\n//--------------------------------------------------------------------------\n// Some would say, most of the magic is done in post! :D\nvec3 PostEffects(vec3 rgb)\n{\n\t// Gamma first...\n\trgb = pow(rgb, vec3(0.45));\n\n\t#define CONTRAST 1.2\n\t#define SATURATION 1.12\n\t#define BRIGHTNESS 1.14\n\treturn mix(vec3(.5), mix(vec3(dot(vec3(.2125, .7154, .0721), rgb*BRIGHTNESS)), rgb*BRIGHTNESS, SATURATION), CONTRAST);\n}\n\n//--------------------------------------------------------------------------\nvoid main(void)\n{\n    vec2 xy = -1.0 + 2.0*gl_FragCoord.xy / resolution.xy;\n\tvec2 s = xy * vec2(resolution.x/resolution.y,1.0);\n\tvec3 camTar;\n\n\t#ifdef STEREO\n\tfloat isCyan = mod(gl_FragCoord.x + mod(gl_FragCoord.y,2.0),2.0);\n\t#endif\n\n\t// Use several forward heights, of decreasing influence with distance from the camera.\n\tfloat h = 0.0;\n\tfloat f = 1.0;\n\tfor (int i = 0; i < 7; i++)\n\t{\n\t\th += Terrain(CameraPath((1.0-f)*.004).xz) * f;\n\t\tf -= .1;\n\t}\n\tcameraPos.xz = CameraPath(0.0).xz;\n\tcamTar.xz\t = CameraPath(.005).xz;\n\tcamTar.y = cameraPos.y = (h*.23)+3.5;\n\t\n\tfloat roll = 0.15*sin(iGlobalTime*.2);\n\tvec3 cw = normalize(camTar-cameraPos);\n\tvec3 cp = vec3(sin(roll), cos(roll),0.0);\n\tvec3 cu = normalize(cross(cw,cp));\n\tvec3 cv = normalize(cross(cu,cw));\n\tvec3 rd = normalize( s.x*cu + s.y*cv + 1.5*cw );\n\n\t#ifdef STEREO\n\tcameraPos += .45*cu*isCyan; // move camera to the right - the rd vector is still good\n\t#endif\n\n\tvec3 col;\n\tfloat distance;\n\tif( !Scene(cameraPos,rd, distance) )\n\t{\n\t\t// Missed scene, now just get the sky value...\n\t\tcol = GetSky(rd);\n\t\tcol = GetClouds(col, rd);\n\t}\n\telse\n\t{\n\t\t// Get world coordinate of landscape...\n\t\tvec3 pos = cameraPos + distance * rd;\n\t\t// Get normal from sampling the high definition height map\n\t\t// Use the distance to sample larger gaps to help stop aliasing...\n\t\tfloat p = min(.3, .0005+.00001 * distance*distance);\n\t\tvec3 nor  \t= vec3(0.0,\t\t    Terrain2(pos.xz), 0.0);\n\t\tvec3 v2\t\t= nor-vec3(p,\t\tTerrain2(pos.xz+vec2(p,0.0)), 0.0);\n\t\tvec3 v3\t\t= nor-vec3(0.0,\t\tTerrain2(pos.xz+vec2(0.0,-p)), -p);\n\t\tnor = cross(v2, v3);\n\t\tnor = normalize(nor);\n\n\t\t// Get the colour using all available data...\n\t\tcol = TerrainColour(pos, nor, distance);\n\t}\n\n\tcol = PostEffects(col);\n\t\n\t#ifdef STEREO\t\n\tcol *= vec3( isCyan, 1.0-isCyan, 1.0-isCyan );\t\n\t#endif\n\t\n\tgl_FragColor=vec4(col.ggr,1.0);\n\tgl_FragColor *= mod(gl_FragCoord.y, 2.0);\n\tgl_FragColor *= vec4(0.9, 1.0, 1.25, 1.0);\n}\n\n//--------------------------------------------------------------------------}" } ] }
{ "_id" : 10141, "created_at" : { "$date" : 1374639573984 }, "image_url" : "/thumbs/10141.png", "modified_at" : { "$date" : 1374639573984 }, "parent" : 9970, "parent_version" : 0, "user" : "ca2d2b8", "versions" : [ { "created_at" : { "$date" : 1374639573984 }, "code" : "#ifdef GL_ES\nprecision mediump float;\n#endif\n\nuniform float time;\nuniform vec2 mouse;\nuniform vec2 resolution;\nfloat iGlobalTime = time;\n\n// Mountains. (C) David Hoskins - 2013\n\n\n// https://www.shadertoy.com/view/4slGD4\n\n// A ray-marched version of my terrain renderer which uses\n// streaming texture normals for speed:-\n// http://www.youtube.com/watch?v=qzkBnCBpQAM\n\n// It was difficult finding a suitable starting point, but I think this one works OK\n\n// It uses binary subdivision to accurately find the height map.\n// Lots of thanks to Iñigo and his noise functions!\n\n// Video of my OpenGL version that \n// http://www.youtube.com/watch?v=qzkBnCBpQAM\n\n// Stereo version code thanks to Croqueteer :)\n//#define STEREO \n\n// Remove the following line to take out the trees.\n#define TREES\n\n#ifdef TREES\nfloat treeLine = 0.0;\nfloat treeCol =100.0;\n#endif\n\nvec3 sunLight  = normalize( vec3(  0.4, 0.4,  0.48 ) );\nvec3 sunColour = vec3(1.0, .9, .83);\nfloat specular = 0.0;\nvec3 cameraPos;\n\n\n// This peturbs the fractal positions for each iteration down...\n// Helps make nice twisted landscapes...\nconst mat2 rotate2D = mat2(1.4623, 1.67231, -1.67231, 1.4623);\n\n// Alternative rotation:-\n// const mat2 rotate2D = mat2(1.2323, 1.999231, -1.999231, 1.22);\n\n//--------------------------------------------------------------------------\n// Noise functions...\nfloat Hash( float n )\n{\n    return fract(sin(n)*43758.5453123);\n}\n\n//--------------------------------------------------------------------------\nfloat Hash(vec2 p)\n{\n\treturn fract(sin(dot(p, vec2(12.9898, 78.233))) * 43758.5453);\n}\n\n//--------------------------------------------------------------------------\nfloat Noise( in vec3 x )\n{\n    vec3 p = floor(x);\n    vec3 f = fract(x);\n    f = f*f*(3.0-2.0*f);\n    float n = p.x + p.y*57.0 + 113.0*p.z;\n    float res = mix(mix(mix( Hash(n+  0.0), Hash(n+  1.0),f.x),\n                        mix( Hash(n+ 57.0), Hash(n+ 58.0),f.x),f.y),\n                    mix(mix( Hash(n+113.0), Hash(n+114.0),f.x),\n                        mix( Hash(n+170.0), Hash(n+171.0),f.x),f.y),f.z);\n    return res;\n}\n//--------------------------------------------------------------------------\nfloat Noise( in vec2 x )\n{\n    vec2 p = floor(x);\n    vec2 f = fract(x);\n    f = f*f*(3.0-2.0*f);\n    float n = p.x + p.y*57.0;\n    float res = mix(mix( Hash(n+  0.0), Hash(n+  1.0),f.x),\n                    mix( Hash(n+ 57.0), Hash(n+ 58.0),f.x),f.y);\n    return res;\n}\n\n//--------------------------------------------------------------------------\nvec2 Noise2( in vec2 x )\n{\n\tvec2 res = vec2(Noise(x), Noise(x+vec2(4101.03, 2310.0)));\n    return res-vec2(.5, .5);\n}\n\n//--------------------------------------------------------------------------\n// iq's derivative noise function...\nvec3 NoiseDerivative( in vec2 x )\n{\n    vec2 p = floor(x);\n    vec2 f = fract(x);\n    vec2 u = f*f*(3.0-2.0*f);\n    float n = p.x + p.y*57.0;\n    float a = Hash(n+  0.0);\n    float b = Hash(n+  1.0);\n    float c = Hash(n+ 57.0);\n    float d = Hash(n+ 58.0);\n\treturn vec3(a+(b-a)*u.x+(c-a)*u.y+(a-b-c+d)*u.x*u.y,\n\t\t\t\t30.0*f*f*(f*(f-2.0)+1.0)*(vec2(b-a,c-a)+(a-b-c+d)*u.yx));\n}\n\n//--------------------------------------------------------------------------\n#ifdef TREES\nfloat Trees(vec2 p)\n{\n\tp *= 5.0;\n\tvec2 rnd = Noise2(p);//vec2(Hash(floor(p.x*4.0)), Hash(floor(p.y*4.0)))*.5;\n\tvec2 v2 = fract(p+rnd)-.5;\n\treturn max(.5-(length(v2)), 0.0) * treeLine*.6;\n}\n#endif\n\n//--------------------------------------------------------------------------\n// Low def version for ray-marching through the height field...\nfloat Terrain( in vec2 p)\n{\n\tvec2 pos = p*0.08;\n\tfloat w = (Noise(pos*.25)*0.75+.15);\n\tw = 36.0 * w * w;\n\tvec2 dxy = vec2(0.0, 0.0);\n\tfloat f = .0;\n\tfor (int i = 0; i < 5; i++)\n\t{\n\t\tvec3 v = NoiseDerivative(pos);\n\t\tdxy += v.yz;\n\t\tf += (w * (v.x) / (1.0 + dot(dxy, dxy))) ;\n\t\tw = -w * 0.37;\t//...Flip negative and positive for variation\n\t\tpos = rotate2D * pos;\n\t}\n\tfloat ff = Noise(pos*.003);\n\t\n\tf += pow(ff, 6.0)*85.-1.0;\n\treturn f;\n}\n\n//--------------------------------------------------------------------------\n// Map to lower resolution for height field mapping for Scene function...\nfloat Map(in vec3 p)\n{\n\tfloat h = Terrain(p.xz);\n\t\t\n\t#ifdef TREES\n\tfloat ff = Noise(p.xz*1.3)*.8;\n\ttreeLine = smoothstep(ff, .1+ff, h) * smoothstep(.5+ff, .4+ff, h);\n\ttreeCol = Trees(p.xz);\n\th += treeCol;\n\t#endif\n\t\n    return p.y - h;\n}\n\n//--------------------------------------------------------------------------\n// High def version only used for grabbing normal information.\nfloat Terrain2( in vec2 p)\n{\n\t// There's some real magic numbers in here! \n\t// The Noise calls add large mountain ranges for more variation over distances...\n\tvec2 pos = p*0.08;\n\tfloat w = (Noise(pos*.25)*0.75+.15);\n\tw = 36.0 * w * w;\n\tvec2 dxy = vec2(0.0, 0.0);\n\tfloat f = .0;\n\tfor (int i = 0; i < 5; i++)\n\t{\n\t\tvec3 v = NoiseDerivative(pos);\n\t\tdxy += v.yz;\n\t\tf += (w * (v.x)  / (1.0 + dot(dxy, dxy)));\n\t\tw =  - w * 0.37;\t//...Flip negative and positive for varition\t   \n\t\tpos = rotate2D * pos;\n\t}\n\tfloat ff = Noise(pos*.003);\n\tf += pow(ff, 6.0)*85.-1.0;\n\t\n\t#ifdef TREES\n\ttreeCol = Trees(p);\n\tf += treeCol;\n\tif (treeCol > 0.0) return f;\n\t#endif\n\t\n\t// That's the last of the low resolution, now go down further for the Normal data...\n\tfor (int i = 0; i < 6; i++)\n\t{\n\t\tvec3 v = NoiseDerivative(pos);\n\t\tdxy += v.yz;\n\t\tf += (w * (v.x) / (1.0 + dot(dxy, dxy)));\n\t\tw =  - w * 0.37;\n\t\tpos = rotate2D * pos;\n\t}\n\t\n\t\n\treturn f;\n}\n\n//--------------------------------------------------------------------------\nfloat FractalNoise(in vec2 xy)\n{\n\tfloat w = .65;\n\tfloat f = 0.0;\n\n\tfor (int i = 0; i < 4; i++)\n\t{\n\t\tf += Noise(xy) * w;\n\t\tw *= 0.5;\n\t\txy *= 2.3;\n\t}\n\treturn f;\n}\n\n//--------------------------------------------------------------------------\n// Simply Perlin clouds that fade to the horizon...\n// 200 units above the ground...\nvec3 GetClouds(in vec3 sky, in vec3 rd)\n{\n\tif (rd.y < 0.0) return sky;\n\tfloat v = (200.0-cameraPos.y)/rd.y;\n\trd.xz *= v;\n\trd.xz += cameraPos.xz;\n\trd.xz *= .010;\n\tfloat f = (FractalNoise(rd.xz) -.55) * 5.0;\n\t// Uses the ray's y component for horizon fade of fixed colour clouds...\n\tsky = mix(sky, vec3(.55, .55, .52), clamp(f*rd.y-.1, 0.0, 1.0));\n\n\treturn sky;\n}\n\n\n\n//--------------------------------------------------------------------------\n// Grab all sky information for a given ray from camera\nvec3 GetSky(in vec3 rd)\n{\n\tfloat sunAmount = max( dot( rd, sunLight), 0.0 );\n\tfloat v = pow(1.0-max(rd.y,0.0),5.)*.5;\n\tvec3  sky = vec3(v*sunColour.x*0.4+0.18, v*sunColour.y*0.4+0.22, v*sunColour.z*0.4+.4);\n\t// Wide glare effect...\n\tsky = sky + sunColour * pow(sunAmount, 6.5)*.32;\n\t// Actual sun...\n\tsky = sky+ sunColour * min(pow(sunAmount, 1150.0), .3)*.65;\n\treturn sky;\n}\n\n//--------------------------------------------------------------------------\n// Merge mountains into te sky background for correct disappearance...\nvec3 ApplyFog( in vec3  rgb, in float dis, in vec3 dir)\n{\n\tfloat fogAmount = clamp(dis* 0.0000165, 0.0, 1.0);\n\treturn mix( rgb, GetSky(dir), fogAmount );\n}\n\n//--------------------------------------------------------------------------\n// Calculate sun light...\nvoid DoLighting(inout vec3 mat, in vec3 pos, in vec3 normal, in vec3 eyeDir, in float dis)\n{\n\tfloat h = dot(sunLight,normal);\n\tfloat c = max(h, 0.0)+.1;\n\tmat = mat * sunColour * c ;\n\t// Specular...\n\tif (h > 0.0)\n\t{\n\t\tvec3 R = reflect(sunLight, normal);\n\t\tfloat specAmount = pow( max(dot(R, normalize(eyeDir)), 0.0), 3.0)*specular;\n\t\tmat = mix(mat, sunColour, specAmount);\n\t}\n}\n\n//--------------------------------------------------------------------------\n// Hack the height, position, and normal data to create the coloured landscape\nvec3 TerrainColour(vec3 pos, vec3 normal, float dis)\n{\n\tvec3 mat;\n\tspecular = .0;\n\tvec3 dir = normalize(pos-cameraPos);\n\t\n\tvec3 matPos = pos * 2.0;// ... I had change scale halfway though, this lazy multiply allow me to keep the graphic scales I had\n\n\tfloat disSqrd = dis * dis;// Squaring it gives better distance scales.\n\n\tfloat f = clamp(Noise(matPos.xz*.05), 0.0,1.0);//*10.8;\n\tf += Noise(matPos.xz*.1+normal.yz*1.08)*.85;\n\tf *= .55;\n\tvec3 m = mix(vec3(.63*f+.2, .7*f+.1, .7*f+.1), vec3(f*.43+.1, f*.3+.2, f*.35+.1), f*.65);\n\tmat = m*vec3(f*m.x+.36, f*m.y+.30, f*m.z+.28);\n\t// Should have used smoothstep to add colours, but left it using 'if' for sanity...\n\tif (normal.y < .5)\n\t{\n\t\tfloat v = normal.y;\n\t\tfloat c = (.5-normal.y) * 4.0;\n\t\tc = clamp(c*c, 0.1, 1.0);\n\t\tf = Noise(vec2(matPos.x*.09, matPos.z*.095+matPos.yy*0.15));\n\t\tf += Noise(vec2(matPos.x*2.233, matPos.z*2.23))*0.5;\n\t\tmat = mix(mat, vec3(.4*f), c);\n\t\tspecular+=.1;\n\t}\n\n\t// Grass. Use the normal to decide when to plonk grass down...\n\tif (matPos.y < 45.35 && normal.y > .65)\n\t{\n\n\t\tm = vec3(Noise(matPos.xz*.073)*.5+.15, Noise(matPos.xz*.12)*.6+.25, 0.0);\n\t\tm *= (normal.y- 0.75)*.85;\n\t\tmat = mix(mat, m, clamp((normal.y-.65)*1.3 * (45.35-matPos.y)*0.1, 0.0, 1.0));\n\t}\n\t#ifdef TREES\n\tif (treeCol > 0.0)\n\t{\n\t\tmat = vec3(.02+Noise(matPos.xz*5.0)*.03, .05, .0);\n\t\tnormal = normalize(normal+vec3(Noise(matPos.xz*33.0)*1.0-.5, .0, Noise(matPos.xz*33.0)*1.0-.5));\n\t\tspecular = .0;\n\t}\n\t#endif\n\t\n\t// Snow topped mountains...\n\tif (matPos.y > 50.0 && normal.y > .28)\n\t{\n\t\tfloat snow = clamp((matPos.y - 50.0 - Noise(matPos.xz * .1)*28.0) * 0.035, 0.0, 1.0);\n\t\tmat = mix(mat, vec3(.7,.7,.8), snow);\n\t\tspecular += snow;\n\t}\n\t// Beach effect...\n\tif (matPos.y < 1.45)\n\t{\n\t\tif (normal.y > .4)\n\t\t{\n\t\t\tf = Noise(matPos.xz * .084)*1.5;\n\t\t\tf = clamp((1.45-f-matPos.y) * 1.34, 0.0, .67);\n\t\t\tfloat t = (normal.y-.4);\n\t\t\tt = (t*t);\n\t\t\tmat = mix(mat, vec3(.09+t, .07+t, .03+t), f);\n\t\t}\n\t\t// Cheap under water darkening...it's wet after all...\n\t\tif (matPos.y < 0.0)\n\t\t{\n\t\t\tmat *= .5;\n\t\t}\n\t}\n\n\tDoLighting(mat, pos, normal,dir, disSqrd);\n\t\n\t// Do the water...\n\tif (cameraPos.y < 0.0)\n\t{\n\t\t// Can go under water, but current camera doesn't find a place...\n\t\tmat = mix(mat, vec3(0.0, .1, .2), .75); \n\t}else\n\tif (matPos.y < 0.0)\n\t{\n\t\t// Pull back along the ray direction to get water surface point at y = 0.0 ...\n\t\tfloat time = (iGlobalTime)*.03;\n\t\tvec3 watPos = matPos;\n\t\twatPos += -dir * (watPos.y/dir.y);\n\t\t// Make some dodgy waves...\n\t\tfloat tx = cos(watPos.x*.052) *4.5;\n\t\tfloat tz = sin(watPos.z*.072) *4.5;\n\t\tvec2 co = Noise2(vec2(watPos.x*4.7+1.3+tz, watPos.z*4.69+time*35.0-tx));\n\t\tco += Noise2(vec2(watPos.z*8.6+time*13.0-tx, watPos.x*8.712+tz))*.4;\n\t\tvec3 nor = normalize(vec3(co.x, 20.0, co.y));\n\t\tnor = normalize(reflect(dir, nor));//normalize((-2.0*(dot(dir, nor))*nor)+dir);\n\t\t// Mix it in at depth transparancy to give beach cues..\n\t\tmat = mix(mat, GetClouds(GetSky(nor), nor), clamp((watPos.y-matPos.y)*1.1, .4, .66));\n\t\t// Add some extra water glint...\n\t\tfloat sunAmount = max( dot(nor, sunLight), 0.0 );\n\t\tmat = mat + sunColour * pow(sunAmount, 228.5)*.6;\n\t}\n\tmat = ApplyFog(mat, disSqrd, dir);\n\treturn mat;\n}\n\n//--------------------------------------------------------------------------\nfloat BinarySubdivision(in vec3 rO, in vec3 rD, float t, float oldT)\n{\n\t// Home in on the surface by dividing by two and split...\n\tfor (int n = 0; n < 4; n++)\n\t{\n\t\tfloat halfwayT = (oldT + t ) * .5;\n\t\tvec3 p = rO + halfwayT*rD;\n\t\tif (Map(p) < 0.25)\n\t\t{\n\t\t\tt = halfwayT;\n\t\t}else\n\t\t{\n\t\t\toldT = halfwayT;\n\t\t}\n\t}\n\treturn t;\n}\n\n//--------------------------------------------------------------------------\nbool Scene(in vec3 rO, in vec3 rD, out float resT )\n{\n    float t = 1.2;\n\tfloat oldT = 0.0;\n\tfloat delta = 0.0;\n\tfor( int j=0; j<170; j++ )\n\t{\n\t\tif (t > 240.0) return false; // ...Too far\n\t    vec3 p = rO + t*rD;\n        if (p.y > 95.0) return false; // ...Over highest mountain\n\n\t\tfloat h = Map(p); // ...Get this positions height mapping.\n\t\t// Are we inside, and close enough to fudge a hit?...\n\t\tif( h < 0.25)\n\t\t{\n\t\t\t// Yes! So home in on height map...\n\t\t\tresT = BinarySubdivision(rO, rD, t, oldT);\n\t\t\treturn true;\n\t\t}\n\t\t// Delta ray advance - a fudge between the height returned\n\t\t// and the distance already travelled.\n\t\t// It's a really fiddly compromise between speed and accuracy\n\t\t// Too large a step and the tops of ridges get missed.\n\t\tdelta = max(0.01, 0.2*h) + (t*0.0065);\n\t\toldT = t;\n\t\tt += delta;\n\t}\n\n\treturn false;\n}\n\n//--------------------------------------------------------------------------\nvec3 CameraPath( float t )\n{\n\tfloat m = 1.0+(mouse.x)*300.0;\n\tt = (iGlobalTime*1.5+m+350.)*.006 + t;\n    vec2 p = 276.0*vec2( sin(3.5*t), cos(1.5*t) );\n\treturn vec3(140.0-p.x, 0.6, -88.0+p.y);\n}\n\n//--------------------------------------------------------------------------\n// Some would say, most of the magic is done in post! :D\nvec3 PostEffects(vec3 rgb)\n{\n\t// Gamma first...\n\trgb = pow(rgb, vec3(0.45));\n\n\t#define CONTRAST 1.2\n\t#define SATURATION 1.12\n\t#define BRIGHTNESS 1.14\n\treturn mix(vec3(.5), mix(vec3(dot(vec3(.2125, .7154, .0721), rgb*BRIGHTNESS)), rgb*BRIGHTNESS, SATURATION), CONTRAST);\n}\n\n//--------------------------------------------------------------------------\nvoid main(void)\n{\n    vec2 xy = -1.0 + 2.0*gl_FragCoord.xy / resolution.xy;\n\tvec2 s = xy * vec2(resolution.x/resolution.y,1.0);\n\tvec3 camTar;\n\n\t#ifdef STEREO\n\tfloat isCyan = mod(gl_FragCoord.x + mod(gl_FragCoord.y,2.0),2.0);\n\t#endif\n\n\t// Use several forward heights, of decreasing influence with distance from the camera.\n\tfloat h = 0.0;\n\tfloat f = 1.0;\n\tfor (int i = 0; i < 7; i++)\n\t{\n\t\th += Terrain(CameraPath((1.0-f)*.004).xz) * f;\n\t\tf -= .1;\n\t}\n\tcameraPos.xz = CameraPath(0.0).xz;\n\tcamTar.xz\t = CameraPath(.005).xz;\n\tcamTar.y = cameraPos.y = (h*.23)+3.5;\n\t\n\tfloat roll = 0.15*sin(iGlobalTime*.2);\n\tvec3 cw = normalize(camTar-cameraPos);\n\tvec3 cp = vec3(sin(roll), cos(roll),0.0);\n\tvec3 cu = normalize(cross(cw,cp));\n\tvec3 cv = normalize(cross(cu,cw));\n\tvec3 rd = normalize( s.x*cu + s.y*cv + 1.5*cw );\n\n\t#ifdef STEREO\n\tcameraPos += .45*cu*isCyan; // move camera to the right - the rd vector is still good\n\t#endif\n\n\tvec3 col;\n\tfloat distance;\n\tif( !Scene(cameraPos,rd, distance) )\n\t{\n\t\t// Missed scene, now just get the sky value...\n\t\tcol = GetSky(rd);\n\t\tcol = GetClouds(col, rd);\n\t}\n\telse\n\t{\n\t\t// Get world coordinate of landscape...\n\t\tvec3 pos = cameraPos + distance * rd;\n\t\t// Get normal from sampling the high definition height map\n\t\t// Use the distance to sample larger gaps to help stop aliasing...\n\t\tfloat p = min(.3, .0005+.00001 * distance*distance);\n\t\tvec3 nor  \t= vec3(0.0,\t\t    Terrain2(pos.xz), 0.0);\n\t\tvec3 v2\t\t= nor-vec3(p,\t\tTerrain2(pos.xz+vec2(p,0.0)), 0.0);\n\t\tvec3 v3\t\t= nor-vec3(0.0,\t\tTerrain2(pos.xz+vec2(0.0,-p)), -p);\n\t\tnor = cross(v2, v3);\n\t\tnor = normalize(nor);\n\n\t\t// Get the colour using all available data...\n\t\tcol = TerrainColour(pos, nor, distance);\n\t}\n\n\tcol = PostEffects(col);\n\t\n\t#ifdef STEREO\t\n\tcol *= vec3( isCyan, 1.0-isCyan, 1.0-isCyan );\t\n\t#endif\n\t\n\tgl_FragColor=vec4(col.ggr,1.0);\n\tgl_FragColor *= mod(gl_FragCoord.y, 2.0);\n\tgl_FragColor *= vec4(0.9, 1.0, 1.25, 1.0);\n}\n\n//--------------------------------------------------------------------------}" } ] }
{ "_id" : 10142, "created_at" : { "$date" : 1374666137474 }, "image_url" : "/thumbs/10142.png", "modified_at" : { "$date" : 1374666296584 }, "parent" : 8035, "parent_version" : 0, "user" : "eeb0063", "versions" : [ { "created_at" : { "$date" : 1374666137474 }, "code" : "#ifdef GL_ES\nprecision mediump float;\n#endif\n\nuniform float time;\nuniform vec2 mouse;\nuniform vec2 resolution;\n\nvoid main() {\n    vec4 color;\n    float y = gl_FragCoord.y;\n    float x = gl_FragCoord.x;\n    \n    vec4 white = vec4(1);\n    vec4 red = vec4(0.3, 1, 1, 1);\n    vec4 blue = vec4(0, 0.2, 0.2, 1);\n    vec4 green = vec4(0, 0.3, 0.8, 1);\n\n    float step1 = resolution.y * 0.00;\n    float step2 = resolution.y * 0.50;\n    float step3 = resolution.y * 0.75;\n\n    color = mix(white, red, smoothstep(step1, step2, y - 0.3*x));\n    color = mix(color, blue, smoothstep(step2, step3, y - 0.3*x));\n    color = mix(color, green, smoothstep(step3, resolution.y, y - 0.3*x));\n\t\n    gl_FragColor = color;\n}\n" }, { "created_at" : { "$date" : 1374666296584 }, "code" : "#ifdef GL_ES\nprecision mediump float;\n#endif\n\nuniform float time;\nuniform vec2 mouse;\nuniform vec2 resolution;\n\nvoid main() {\n    vec4 color;\n    float y = gl_FragCoord.y;\n    float x = gl_FragCoord.x;\n    \n    vec4 white = vec4(1);\n    vec4 red = vec4(0.8, 0.3, 0.3, 1);\n    vec4 blue = vec4(0.3, 0.6, 1, 1);\n    vec4 green = vec4(0.3, 1, 0.3, 1);\n\n    float step1 = resolution.y * 0.00;\n    float step2 = resolution.y * 0.50;\n    float step3 = resolution.y * 0.75;\n\n    color = mix(white, red, smoothstep(step1, step2, y - 0.3*x));\n    color = mix(color, blue, smoothstep(step2, step3, y - 0.3*x));\n    color = mix(color, green, smoothstep(step3, resolution.y, y - 0.3*x));\n\t\n    gl_FragColor = color;\n}\n" } ] }
{ "_id" : 10143, "created_at" : { "$date" : 1374667807852 }, "image_url" : "/thumbs/10143.png", "modified_at" : { "$date" : 1374667807853 }, "parent" : 10130, "parent_version" : 0, "user" : "abbdc60", "versions" : [ { "created_at" : { "$date" : 1374667807853 }, "code" : "// Playing around with Lissajous curves.\n#ifdef GL_ES\nprecision mediump float;\n#endif\n\nuniform float time;\nuniform vec2 resolution;\n\nconst int num = 100;\n\nvoid main( void ) {\n    float sum = 0.0;\n\t\n    float size = resolution.x / 500.0;\n\t\n    for (int i = 0; i < num; ++i) {\n        vec2 position = resolution / 2.0;\n\tfloat t = (float(i) + time) / 5.0;\n\tfloat c = float(i) * 4.0;\n        position.x += tan(8.0 * t + c) * resolution.x * 0.2;\n        position.y += sin(t) * resolution.y * .8;\n\n        sum += size / length(gl_FragCoord.xy - position);\n    }\n\t\n    gl_FragColor = vec4(sum * 0.1, sum * 0.5, sum, 1);\n}" } ] }`
)

var (
	testEffects = []Effect{
		{
			ID:            1,
			CreatedAt:     time.Now(),
			ModifiedAt:    time.Date(2021, time.January, 2, 0, 0, 0, 0, time.UTC),
			Parent:        0,
			ParentVersion: 0,
			User:          "1",
			Hidden:        false,
		},
		{
			ID:            2,
			CreatedAt:     time.Now(),
			ModifiedAt:    time.Date(2021, time.January, 1, 0, 0, 0, 0, time.UTC),
			Parent:        0,
			ParentVersion: 0,
			User:          "2",
			Hidden:        true,
		},
	}
)
