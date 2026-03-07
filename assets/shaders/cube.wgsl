// Camera uniform: view-projection matrix.
struct CameraUniforms {
    viewProj: mat4x4<f32>,
}

// Per-mesh uniform: model matrix + normal matrix (transpose(inverse(model))).
struct MeshUniforms {
    model: mat4x4<f32>,
    normalMatrix: mat4x4<f32>,
}

// Directional light parameters (each vec3 padded to vec4 for alignment).
struct LightUniforms {
    direction: vec3<f32>,  // normalized light direction (world-space)
    _pad0: f32,
    color: vec3<f32>,      // light RGB intensity
    _pad1: f32,
    ambient: vec3<f32>,    // ambient light RGB
    _pad2: f32,
    cameraPos: vec3<f32>,  // camera world position (for specular)
    _pad3: f32,
}

@group(0) @binding(0) var<uniform> camera: CameraUniforms;
@group(0) @binding(1) var<uniform> mesh: MeshUniforms;
@group(0) @binding(2) var<uniform> light: LightUniforms;

// Albedo texture and sampler (group 1).
// When no texture is provided, a 1x1 white default texture is bound.
@group(1) @binding(0) var albedo_texture: texture_2d<f32>;
@group(1) @binding(1) var albedo_sampler: sampler;

struct VertexInput {
    @location(0) position: vec3<f32>,
    @location(1) color: vec3<f32>,
    @location(2) normal: vec3<f32>,
    @location(3) uv: vec2<f32>,
}

struct VertexOutput {
    @builtin(position) clip_position: vec4<f32>,
    @location(0) color: vec3<f32>,
    @location(1) world_normal: vec3<f32>,
    @location(2) world_pos: vec3<f32>,
    @location(3) uv: vec2<f32>,
}

@vertex
fn vs_main(in: VertexInput) -> VertexOutput {
    var out: VertexOutput;
    let world_pos = mesh.model * vec4<f32>(in.position, 1.0);
    out.clip_position = camera.viewProj * world_pos;
    out.color = in.color;
    out.world_pos = world_pos.xyz;
    // Transform normal by the upper-left 3x3 of the normal matrix.
    out.world_normal = normalize((mesh.normalMatrix * vec4<f32>(in.normal, 0.0)).xyz);
    out.uv = in.uv;
    return out;
}

@fragment
fn fs_main(in: VertexOutput) -> @location(0) vec4<f32> {
    let N = normalize(in.world_normal);
    let L = normalize(-light.direction);

    // Albedo from texture * vertex color (white texture = vertex color only).
    let tex_color = textureSample(albedo_texture, albedo_sampler, in.uv).rgb;
    let albedo = tex_color * in.color;

    // Diffuse (Lambertian).
    let NdotL = max(dot(N, L), 0.0);
    let diffuse = light.color * NdotL;

    // Specular (Blinn-Phong).
    let V = normalize(light.cameraPos - in.world_pos);
    let H = normalize(L + V);
    let NdotH = max(dot(N, H), 0.0);
    let shininess = 32.0;
    let specular = light.color * pow(NdotH, shininess);

    let result = (light.ambient + diffuse) * albedo + specular * 0.3;
    return vec4<f32>(result, 1.0);
}
