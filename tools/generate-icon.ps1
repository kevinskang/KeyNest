Add-Type -AssemblyName System.Drawing

$size = 512
$bmp = New-Object System.Drawing.Bitmap($size, $size, [System.Drawing.Imaging.PixelFormat]::Format32bppArgb)
$g = [System.Drawing.Graphics]::FromImage($bmp)
$g.SmoothingMode    = [System.Drawing.Drawing2D.SmoothingMode]::AntiAlias
$g.CompositingMode  = [System.Drawing.Drawing2D.CompositingMode]::SourceOver
$g.Clear([System.Drawing.Color]::Transparent)

# ---- Background (rounded rectangle, dark navy gradient) ----
$bgRadius = 80
$bgPath = New-Object System.Drawing.Drawing2D.GraphicsPath
$bgPath.AddArc(0,           0,           $bgRadius*2, $bgRadius*2, 180, 90)
$bgPath.AddArc($size-$bgRadius*2, 0,           $bgRadius*2, $bgRadius*2, 270, 90)
$bgPath.AddArc($size-$bgRadius*2, $size-$bgRadius*2, $bgRadius*2, $bgRadius*2, 0,   90)
$bgPath.AddArc(0,           $size-$bgRadius*2, $bgRadius*2, $bgRadius*2, 90,  90)
$bgPath.CloseFigure()

$bgBrush = New-Object System.Drawing.Drawing2D.LinearGradientBrush(
    (New-Object System.Drawing.RectangleF(0, 0, $size, $size)),
    [System.Drawing.Color]::FromArgb(255, 20, 48, 90),
    [System.Drawing.Color]::FromArgb(255, 38, 82, 150),
    [System.Drawing.Drawing2D.LinearGradientMode]::ForwardDiagonal
)
$g.FillPath($bgBrush, $bgPath)

# ---- Key shape (gold) ----
$keyColor = [System.Drawing.Color]::FromArgb(255, 245, 196, 40)
$keyBrush = New-Object System.Drawing.SolidBrush($keyColor)

# Shaft  (horizontal bar, right half of icon)
$shaftLeft   = 248
$shaftTop    = 234
$shaftWidth  = 206
$shaftHeight = 44
$g.FillRectangle($keyBrush, $shaftLeft, $shaftTop, $shaftWidth, $shaftHeight)

# Tooth 1 (longer)
$g.FillRectangle($keyBrush, 310, ($shaftTop + $shaftHeight), 38, 68)
# Tooth 2 (shorter)
$g.FillRectangle($keyBrush, 380, ($shaftTop + $shaftHeight), 38, 50)

# Key ring / bow (donut shape) — drawn last so hole shows through
# Outer circle: center (175, 256), radius 108
# Inner circle: center (175, 256), radius 48  → FillMode.Alternate creates donut
$ringPath = New-Object System.Drawing.Drawing2D.GraphicsPath
$ringPath.FillMode = [System.Drawing.Drawing2D.FillMode]::Alternate
$ringPath.AddEllipse(67, 148, 216, 216)    # outer  cx=175 cy=256 r=108
$ringPath.AddEllipse(127, 208, 96, 96)    # hole   cx=175 cy=256 r=48
$g.FillPath($keyBrush, $ringPath)

# ---- Save ----
$scriptDir  = Split-Path -Parent $MyInvocation.MyCommand.Path
$outputPath = Join-Path $scriptDir "..\build\appicon.png"
$outputPath = [System.IO.Path]::GetFullPath($outputPath)
$bmp.Save($outputPath, [System.Drawing.Imaging.ImageFormat]::Png)
$g.Dispose()
$bmp.Dispose()

Write-Host "Key icon saved to: $outputPath"
