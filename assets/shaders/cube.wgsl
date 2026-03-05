struct VertexInput {
    @location(0) position: vec3<f32>,
    @location(1) color: vec3<f32>,
}

struct VertexOutput {
    @builtin(position) clip_position: vec4<f32>,
    @location(0) color: vec3<f32>,
}

// MVP = Projection * View; camera at (3,2,5) looking at origin, fov=60 deg.
// Columns of the column-major mat4x4<f32> — substituted at runtime by pkg/gpu
// from the window aspect ratio. Uniform buffers replace this in Phase 14.
const mvp = mat4x4<f32>(
    vec4<f32>(%f, %f, %f, %f),
    vec4<f32>(%f, %f, %f, %f),
    vec4<f32>(%f, %f, %f, %f),
    vec4<f32>(%f, %f, %f, %f),
);

@vertex
fn vs_main(in: VertexInput) -> VertexOutput {
    var out: VertexOutput;
    out.clip_position = mvp * vec4<f32>(in.position, 1.0);
    out.color = in.color;
    return out;
}

@fragment
fn fs_main(in: VertexOutput) -> @location(0) vec4<f32> {
    return vec4<f32>(in.color, 1.0);
}
