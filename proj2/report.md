# Project: Twitter Feed (Modeled as a Bounded-Buffer Problem)

### Project Description - 

This project is an implementation of an image editor. Images are transformed based on a sequence of effects that are to be applied to each image along with the name of the image to be transformed and the name of the transformed image to be saved. These transformations (effects) are carried by applying specific kernels through a convolution operation over the entire image. For all effects, a 3x3 kernel is convolved over patches (iterating over the entire image). The editor processes image tasks specified in three different ways i.e. sequential, bulk synchronous parallel (bsp) and pipeline parallel.

### Important System Components - 

1. Data (Input):
    - **I/O takes places in the `proj2/data` directory. If there is no `data` directory, you can create one using `mkdir data` while inside `proj2`.**
    - Each editor task is provided within the directory `proj2/data/effects.txt` file. An example of an editor task is shown below (multiple tasks can be specified by placing one on each line of the `effects.txt` file and are in the `JSON format` - 

        ```JSON
        {"inPath": "IMG_2029.png", "outPath": "IMG_2029_Out.png", "effects": ["G","E","S"]}
        ``` 
    
    - Input images must be provided within the a directory inside `proj2/data/in` directory. If one does not exist, you can create it from within `proj2/data` using `mkdir in`. Similarly for the directory. For example, to process an image, we can create  `proj2/data/in/small` and place the image inside this directory.  
    - Input images are provided in the `PNG` format.


2. Data (Output):
    - **I/O takes places in the `proj2/data` directory. If there is no `data` directory, you can create one using `mkdir data` while inside `proj2`.**
    - After processing images, the transformed images are saved within the `proj2/data/out` directory. If one does not exist, you can create it from from within `proj2/data` using `mkdir out`. 
    - Processed images are saved in the `PNG` format.

3. Effects:
    - The following effects can be applied (the appropriate identifier to be used is specified in parenthesis for each effect) - 
        - Grayscale (G) - Each pixel's color channels are computed by averaging over the each the original pixel's color channel values.
        - Sharpen (S) - The following kernel is used as part of the convolution operation applied at every pixel using the convolution operation - 
            
            $
            \begin{bmatrix}
                0 & -1 & 0\\
                -1 & 5 & -1\\
                0 & -1 & 0
            \end{bmatrix}
            $
        
        - Blur (B) - The following kernel is used as part of the convolution operation applied at every pixel using the convolution operation - 
            
            $
            \frac{1}{9}\begin{bmatrix}
                1 & 1 & 1\\
                1 & 1 & 1\\
                1 & 1 & 1
            \end{bmatrix}
            $
        
        - Edge Detection (E) - The following kernel is used as part of the convolution operation applied at every pixel using the convolution operation - 
            
            $
            \begin{bmatrix}
                -1 & -1 & -1\\
                -1 & 8 & -1\\
                -1 & -1 & -1
            \end{bmatrix}
            $

    - **Note that these identifiers for each effect are case sensitive.**
    - Convolution is perform on a 2D image grid by convolving the kernel across the image. The convolution operation can be thought of as a sliding window computation over the entire image. The computation being performed is the frobenius inner product which is the sum over element wise products between the kernel and patch of image overlapping with the kernel. **Zero padding** is used at the edges of the image. The convolution operation can be seen below - 

        $
        y[m,n] = x[m,n] \ast h[m,n] = \sum_{j=-\infty}^{\infty} \sum_{i=-\infty}^{\infty} x[i,j] h[m - i, n - j]
        $

4. Run Modes:
    The following modes can be used with the identifier specified in parenthesis (for sequential mode, no mode is specified when running the program).
    - sequential - If no mode is specified, then the sequential version is run.
    - pipeline (pipeline) - This runs the parallel mode pipeline, which utilizes the fan-out fan-in parallelism scheme. The number of threads must be specified in this mode. Each thread spawned will work on an image tasks (all of the effects for the image task) by spawning the same number of threads for each image task, each working on a portion of the image for that effect.
    - bulk synchronous parallel (bsp) - This runs the parallel mode bsp which utilizes a bulk synchronous parallel scheme in which the number of threads (required argument) specified are spawned to work on a single effect of a single image at a time, where thread works on a patch of the image for that effect.



### Running the Program - 

The editor can be run in the following way - 

1. Sequential run - 

```console
foo@bar:~$ go run path/to/editor.go <directory>
```

2. Parallel runs - 

    1. Pipeline - 

    ```console
    foo@bar:~$ go run path/to/editor.go <image directory> pipeline <number of threads to be spawned>
    ```

    2. BSP - 

    ```console
    foo@bar:~$ go run path/to/editor.go <image directory> bsp <number of threads to be spawned>
    ```

3. Multiple Input Image Directories - 

If there a images in multiple directories within `proj2/data/in`, for example, if there was the directories `small` and `big`, we can chain directories to process using `+` - 

```console
foo@bar:~$ go run path/to/editor.go small+big pipeline <number of threads to be spawned>
```


### Benchmarking the Program - 

The program can be benchmarked using the following command - 

```console
foo@bar:~$ sbatch benchmark-proj2.sh
```

This must be run within the `benchmark directory`. Make sure to create the `slurm/out` directory inside `benchmark` directory and check the `.stdout` file for the outputs and the timings. 

### Benchmarking the Program - 

Benchmarking is done by using the `benchmark/benchmark.go`. To benchmark the application, use the following command when within the `benchmark` directory - 

```console
foo@bar:~$ sbatch benchmark.sh
```

The graph of speedups obtained can be seen below - 

1. Pipeline mode - 

    ![benchmarking_pipeline](./benchmark/pipeline-speedup.png)

2. BSP mode - 

    ![benchmarking_bsp](./benchmark/bsp-speedup.png)


The graphs will be created within the `benchmark` directory. The computation of the speedups along with the storing of each of the benchmarking timings and the plotting of the stored data happens by using `benchmark_graph.py` which is called from within `benchmark-proj2.sh` (both reside in the `benchmark` directory).


The following observations can be made from the graph - 

1. 


### Questions About Implementation - 

1. 